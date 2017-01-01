package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// The base divisor to use for crits.
const baseCritDiv = 7

// A 'character sheet' for an actor. Basically, all the attributes for an actor
// that are required to make calculations and game decisions live here.
type Sheet interface {
	Objgetter

	// Stats
	Str() int
	Agi() int
	Vit() int
	Mnd() int

	// Skills.
	Melee() int
	Evasion() int

	// Attack and defense.
	Attack() Attack
	Defense() Defense

	// Vitals.
	HP() int
	MaxHP() int
	MP() int
	MaxMP() int

	// Hurt me.
	Hurt(dmg int)
	// Heal me.
	Heal(amt int)
	// Touch me.
	// Feel me.
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet struct {
	Trait
	str int
	agi int
	vit int
	mnd int

	hp int
	mp int
}

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{
		Trait: Trait{obj: obj},
		str:   3,
		agi:   4,
		vit:   4,
		mnd:   3,
	}
	ps.hp = ps.MaxHP()
	ps.mp = ps.MaxMP()
	return ps
}

func (s *PlayerSheet) Str() int {
	return s.str
}

func (s *PlayerSheet) Agi() int {
	return s.agi
}

func (s *PlayerSheet) Vit() int {
	return s.vit
}

func (s *PlayerSheet) Mnd() int {
	return s.mnd
}

func (p *PlayerSheet) Melee() int {
	return p.Agi()
}

func (p *PlayerSheet) Evasion() int {
	return p.Agi()
}

func (p *PlayerSheet) HP() int {
	return p.hp
}

func (p *PlayerSheet) MP() int {
	return p.mp
}

func (p *PlayerSheet) MaxHP() int {
	return 10 * (1 + p.Vit())
}

func (p *PlayerSheet) MaxMP() int {
	return 10 * (1 + p.Mnd())
}

func (p *PlayerSheet) Hurt(dmg int) {
	p.hp -= dmg
	checkDeath(p)
}

func (p *PlayerSheet) Heal(amt int) {
	p.hp = math.Min(p.hp+amt, p.MaxHP())
}

func (p *PlayerSheet) Attack() Attack {
	melee := p.obj.Equipper.Body().Melee() + p.Melee()

	weap := p.weapon()
	str := p.Str()

	bonusSides := math.Min(math.Abs(str), weap.Equipment.Weight) * math.Sgn(str)
	return Attack{
		Melee:   melee,
		Damroll: weap.Equipment.Damroll.Add(0, bonusSides),
		CritDiv: weap.Equipment.Weight + baseCritDiv,
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
				Damroll: NewDice(1, p.Str()+1),
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
	return Defense{
		Evasion:  evasion,
		ProtDice: dice,
	}
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	Trait

	str int
	agi int
	vit int
	mnd int

	hp    int
	mp    int
	maxhp int
	maxmp int

	melee   int
	evasion int
	// Basically weapon weight.
	critdivmod int

	protroll Dice
	damroll  Dice
}

// Given a copy of a MonsterSheet literal, this will return a function that will bind
// the owner of the sheet to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewMonsterSheet(sheetspec MonsterSheet) func(*Obj) Sheet {
	return func(o *Obj) Sheet {
		// Copy sheet.
		sheet := sheetspec
		sheet.obj = o
		sheet.hp = sheet.maxhp
		sheet.mp = sheet.maxmp
		return &sheet
	}
}

func (s *MonsterSheet) Str() int {
	return s.str
}

func (s *MonsterSheet) Agi() int {
	return s.agi
}

func (s *MonsterSheet) Vit() int {
	return s.vit
}

func (s *MonsterSheet) Mnd() int {
	return s.mnd
}

func (m *MonsterSheet) Melee() int {
	return m.melee
}

func (m *MonsterSheet) Evasion() int {
	return m.evasion
}

func (m *MonsterSheet) HP() int {
	return m.hp
}

func (m *MonsterSheet) MP() int {
	return m.mp
}

func (m *MonsterSheet) MaxHP() int {
	return m.maxhp
}

func (m *MonsterSheet) MaxMP() int {
	return m.maxmp
}

func (m *MonsterSheet) Hurt(dmg int) {
	m.hp -= dmg
	checkDeath(m)
}

func (m *MonsterSheet) Heal(amt int) {
	m.hp = math.Min(m.hp+amt, m.MaxHP())
}

func (m *MonsterSheet) Attack() Attack {
	return Attack{
		Melee:   m.melee,
		Damroll: m.damroll,
		CritDiv: m.critdivmod + baseCritDiv,
	}
}

func (m *MonsterSheet) Defense() Defense {
	return Defense{
		Evasion:  m.evasion,
		ProtDice: []Dice{m.protroll},
	}
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

// Details about an actor's melee attack, before the melee roll is applied --
// i.e. what melee bonus + damage should be done if no crits happen? What's the
// base divisor to use for calculating # of crits?
type Attack struct {
	Melee   int
	Damroll Dice
	CritDiv int
}

// Roll damage for this attack, given that `crits` crits were rolled.
func (atk Attack) RollDamage(crits int) int {
	return atk.Damroll.Add(crits, 0).Roll()
}

// Details about an actor's defense, before the evasion roll is applied. i.e.
// what evasion bonus should be added and what protection dice should be rolled
// when attacked?
type Defense struct {
	Evasion  int
	ProtDice []Dice
}

func (def Defense) RollProt() int {
	dice := def.ProtDice
	sum := 0

	for _, d := range dice {
		sum += d.Roll()
	}

	return sum
}
