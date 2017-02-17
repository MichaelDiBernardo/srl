package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A 'character sheet' for an actor. Basically, all the attributes for an actor
// that are required to make calculations and game decisions live here.
type Sheet interface {
	Objgetter

	// Stats
	Stat(stat StatName) int
	SetStat(stat StatName, amt int)

	UnmodStat(stat StatName) int
	StatMod(stat StatName) int
	SetStatMod(stat StatName, amt int)

	// Skills.
	Melee() int
	Evasion() int

	// Attack and defense.
	Attack() Attack
	Defense() Defense

	// Vitals.
	HP() int
	setHP(hp int)
	MaxHP() int

	MP() int
	setMP(mp int)
	MaxMP() int

	// Regen factor. Normal healing is 1; 2 is twice as fast, etc.
	// 0 means no regen.
	Regen() int

	// How stunned is this actor?
	Stun() StunLevel

	// Sight radius.
	Sight() int

	// Remove 'dmg' hitpoints. Will kill the actor if HP falls <= 0.
	Hurt(dmg int)
	// Heal 'amt' hp; actor's hp will not exceed maxhp.
	Heal(amt int)
	// Set stun level. Effects are documented with the stun level definitions
	// in this module.
	SetStun(lvl StunLevel)

	Speed() int
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet struct {
	Trait
	stats *stats

	sight int
	speed int

	hp    int
	mp    int
	regen int
	stun  StunLevel
}

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{
			stats: statlist{
				Str: 2,
				Agi: 4,
				Vit: 4,
				Mnd: 3,
			},
		},
		speed: 2,
		regen: 1,
		sight: FOVRadius,
	}
	ps.hp = ps.MaxHP()
	ps.mp = ps.MaxMP()
	return ps
}

func (p *PlayerSheet) Stat(stat StatName) int {
	return p.stats.stat(stat)
}

func (p *PlayerSheet) SetStat(stat StatName, amt int) {
	p.stats.set(stat, amt)
	// TODO: Modify skills if agi or mnd.
}

func (p *PlayerSheet) UnmodStat(stat StatName) int {
	return p.stats.unmodstat(stat)
}

func (p *PlayerSheet) StatMod(stat StatName) int {
	return p.stats.mod(stat)
}

func (p *PlayerSheet) SetStatMod(stat StatName, amt int) {
	p.stats.setmod(stat, amt)
	// TODO: Modify skills if agi or mnd.
}

func (p *PlayerSheet) Melee() int {
	return p.stats.stat(Agi)
}

func (p *PlayerSheet) Evasion() int {
	return p.stats.stat(Agi)
}

func (p *PlayerSheet) Speed() int {
	return p.speed
}

func (p *PlayerSheet) HP() int {
	return p.hp
}

func (p *PlayerSheet) setHP(hp int) {
	p.hp = hp
}

func (p *PlayerSheet) MaxHP() int {
	return 10 * (1 + p.stats.stat(Vit))
}

func (p *PlayerSheet) MP() int {
	return p.mp
}

func (p *PlayerSheet) setMP(mp int) {
	p.mp = mp
}

func (p *PlayerSheet) MaxMP() int {
	return 10 * (1 + p.stats.stat(Mnd))
}

func (p *PlayerSheet) Regen() int {
	return p.regen
}

func (p *PlayerSheet) Sight() int {
	return p.sight
}

func (p *PlayerSheet) Hurt(dmg int) {
	hurt(p, dmg)
}

func (p *PlayerSheet) Heal(amt int) {
	heal(p, amt)
}

func (p *PlayerSheet) Stun() StunLevel {
	return p.stun
}

func (p *PlayerSheet) SetStun(level StunLevel) {
	changestun(p, level)
	p.stun = level
}

func (p *PlayerSheet) Attack() Attack {
	melee := p.obj.Equipper.Body().Melee() + p.Melee()

	weap := p.weapon()
	equip := weap.Equipment

	str := p.stats.stat(Str)
	bonusSides := math.Min(math.Abs(str), weap.Equipment.Weight) * math.Sgn(str)

	return Attack{
		Melee:   melee,
		Damroll: equip.Damroll.Add(0, bonusSides),
		CritDiv: equip.Weight + baseCritDiv,
		Effects: equip.Effects,
	}
}

func (p *PlayerSheet) weapon() *Obj {
	weap := p.obj.Equipper.Body().Weapon()
	if weap != nil {
		return weap
	}
	return p.fist()
}

func (p *PlayerSheet) fist() *Obj {
	return p.obj.Game.NewObj(&Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecFist,
		Name:    "FIST",
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Damroll: NewDice(1, p.stats.stat(Str)+1),
				Melee:   0,
				Weight:  0,
				Slot:    SlotHand,
			}),
		},
	})
}

func (p *PlayerSheet) Defense() Defense {
	body := p.obj.Equipper.Body()
	evasion := body.Evasion() + p.Evasion()
	dice := body.ProtDice()
	effects := body.ArmorEffects()

	return Defense{
		Evasion:  evasion,
		ProtDice: dice,
		Effects:  effects,
	}
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	Trait

	stats *stats

	speed int
	sight int

	hp    int
	mp    int
	maxhp int
	maxmp int

	regen int
	stun  StunLevel

	melee   int
	evasion int
	// Basically weapon weight.
	critdivmod int

	protroll Dice
	damroll  Dice

	// Innate attack + defense effects for this monster.
	atkeffects Effects
	defeffects Effects
}

// Given a copy of a MonsterSheet literal, this will return a function that will bind
// the owner of the sheet to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewMonsterSheet(sheetspec MonsterSheet) func(*Obj) Sheet {
	return func(o *Obj) Sheet {
		// Copy sheet.
		sheet := sheetspec

		// Copy stats.
		// TODO: Just write Copy() on sheet, god.
		newstats := &stats{}
		if sheetspec.stats != nil {
			*newstats = *(sheetspec.stats)
		}
		sheet.stats = newstats

		sheet.obj = o
		sheet.hp = sheet.maxhp
		sheet.mp = sheet.maxmp
		if sheet.regen == 0 {
			sheet.regen = 1
		}
		if sheet.sight == 0 {
			sheet.sight = FOVRadius
		}
		return &sheet
	}
}

func (m *MonsterSheet) Stat(stat StatName) int {
	return m.stats.stat(stat)
}

func (m *MonsterSheet) SetStat(stat StatName, amt int) {
	m.stats.set(stat, amt)
}

func (m *MonsterSheet) UnmodStat(stat StatName) int {
	return m.stats.unmodstat(stat)
}

func (m *MonsterSheet) StatMod(stat StatName) int {
	return m.stats.mod(stat)
}

func (m *MonsterSheet) SetStatMod(stat StatName, amt int) {
	m.stats.setmod(stat, amt)
}

func (m *MonsterSheet) Melee() int {
	return m.melee
}

func (m *MonsterSheet) Evasion() int {
	return m.evasion
}

func (m *MonsterSheet) Speed() int {
	return m.speed
}

func (m *MonsterSheet) HP() int {
	return m.hp
}

func (m *MonsterSheet) setHP(hp int) {
	m.hp = hp
}

func (m *MonsterSheet) MaxHP() int {
	return m.maxhp
}

func (m *MonsterSheet) MP() int {
	return m.mp
}

func (m *MonsterSheet) setMP(mp int) {
	m.mp = mp
}

func (m *MonsterSheet) MaxMP() int {
	return m.maxmp
}

func (m *MonsterSheet) Regen() int {
	return m.regen
}

func (m *MonsterSheet) Sight() int {
	return m.sight
}

func (m *MonsterSheet) Hurt(dmg int) {
	hurt(m, dmg)
}

func (m *MonsterSheet) Heal(amt int) {
	heal(m, amt)
}

func (m *MonsterSheet) Stun() StunLevel {
	return m.stun
}

func (m *MonsterSheet) SetStun(level StunLevel) {
	changestun(m, level)
	m.stun = level
}

func (m *MonsterSheet) Attack() Attack {
	return Attack{
		Melee:   m.melee,
		Damroll: m.damroll,
		CritDiv: m.critdivmod + baseCritDiv,
		Effects: m.atkeffects,
	}
}

func (m *MonsterSheet) Defense() Defense {
	return Defense{
		Evasion:  m.evasion,
		ProtDice: []Dice{m.protroll},
		Effects:  m.defeffects,
	}
}

// Maintains stats and modifications for a sheet.
type stats struct {
	stats statlist
	mods  statlist
}

// An array of statistics.
type statlist [NumStats]int

// Get the given stat combined with any active modifiersto it.
func (s *stats) stat(stat StatName) int {
	return s.stats[stat] + s.mods[stat]
}

// Get the unmodified stat.
func (s *stats) unmodstat(stat StatName) int {
	return s.stats[stat]
}

// Set the given stat to 'amt'.
func (s *stats) set(stat StatName, amt int) {
	s.stats[stat] = amt
}

// Get the modifier (mod) for this stat.
func (s *stats) mod(stat StatName) int {
	return s.mods[stat]
}

// Set a modifier for this stat. Can be +ve or -ve.
func (s *stats) setmod(stat StatName, amt int) {
	s.mods[stat] = amt
}

// Index into stats arrays.
type StatName uint

// The actual stats our actors have.
const (
	Str StatName = iota
	Agi
	Vit
	Mnd
	NumStats
)

// Details about an actor's melee attack, before the melee roll is applied --
// i.e. what melee bonus + damage should be done if no crits happen? What's the
// base divisor to use for calculating # of crits?
type Attack struct {
	Melee   int
	Damroll Dice
	CritDiv int
	Effects Effects
}

// Roll damage for this attack, given that `crits` crits were rolled.
func (atk Attack) RollDamage(extradice int) int {
	return atk.Damroll.Add(extradice, 0).Roll()
}

// Details about an actor's defense, before the evasion roll is applied. i.e.
// what evasion bonus should be added and what protection dice should be rolled
// when attacked?
type Defense struct {
	Evasion  int
	ProtDice []Dice
	Effects  Effects
}

func (def Defense) RollProt() int {
	dice := def.ProtDice
	sum := 0

	for _, d := range dice {
		sum += d.Roll()
	}

	return sum
}

// Shared functions across sheets.
func heal(s Sheet, amt int) {
	s.setHP(math.Min(s.HP()+amt, s.MaxHP()))
}

func hurt(s Sheet, dmg int) {
	s.setHP(s.HP() - dmg)
	checkDeath(s)
}

func changestun(s Sheet, newstun StunLevel) {
	oldstun := s.Stun()
	if oldstun >= newstun {
		return
	}

	msg := s.Obj().Spec.Name + " is "

	switch newstun {
	case Stunned:
		msg += "stunned."
	case MoreStunned:
		msg += "more stunned."
	case KnockedOut:
		msg += "knocked out!"
	}

	s.Obj().Game.Events.Message(msg)
}

func checkDeath(s Sheet) {
	if s.HP() > 0 {
		return
	}

	obj := s.Obj()
	game := obj.Game

	game.Events.Message(fmt.Sprintf("%s fell.", obj.Spec.Name))
	game.Kill(obj)
}

type StunLevel uint

// Stun status definitions.
const (
	// Not stunned, no effect.
	NotStunned StunLevel = iota
	// -2 to all skills
	Stunned
	// -4 to all skills
	MoreStunned
	// knocked out, cannot move.
	KnockedOut
)

// The base divisor to use for crits.
const baseCritDiv = 7
