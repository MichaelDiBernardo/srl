package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A 'character sheet' for an actor. Basically, all the attributes for an actor
// that are required to make calculations and game decisions live here.
type Sheet interface {
	Objgetter

	// Get the total value for statistic 'stat', including mods.
	Stat(stat StatName) int
	// Set the base value for statistic 'stat'.
	SetStat(stat StatName, amt int)

	// Get the unmodified (base) value of statistic 'stat'.
	UnmodStat(stat StatName) int
	// Get the mod for statistic 'stat'.
	StatMod(stat StatName) int
	// Change the mod for 'stat' by 'diff'.
	ChangeStatMod(stat StatName, diff int)

	// Get the total value for skill 'skill', including mods.
	Skill(skill SkillName) int
	// Set the base value for skill 'skill'.
	SetSkill(skill SkillName, amt int)

	// Get the unmodified (base) value of skill 'skill'.
	UnmodSkill(skill SkillName) int
	// Get the mod for skill 'skill'.
	SkillMod(skill SkillName) int
	// Change the mod for skill 'skill'.
	ChangeSkillMod(skill SkillName, diff int)

	// Get information about this actor's melee attack capability.
	Attack() Attack
	// Get information about this actor's defensive capability.
	Defense() Defense

	// Get this actor's current HP.
	HP() int
	// Set this actor's current HP. This will ignore MaxHP constraints.
	setHP(hp int)
	// Get this actor's maximum HP.
	MaxHP() int

	// Does this actor still have > 0 HP?
	Dead() bool

	// Remove 'dmg' hitpoints. Will kill the actor if HP falls <= 0.
	Hurt(dmg int)
	// Heal 'amt' hp; actor's hp will not exceed maxhp.
	Heal(amt int)

	// Get this actor's current MP.
	MP() int
	// Set this actor's current MP. This will ignore MaxMP constraints.
	setMP(mp int)
	// Get this actor's current MP.
	MaxMP() int

	// Remove 'amt' mp. Will clamp MP to 0.
	HurtMP(amt int)
	// Add 'amt' mp; actor's mp will not exceed maxmp.
	HealMP(amt int)

	// Regen factor. Normal healing is 1; 2 is twice as fast, etc.
	// 0 means no regen.
	Regen() int

	// How stunned is this actor? See StunLevel definition to see what each
	// level means.
	Stun() StunLevel
	// Set stun level. Effects are documented with the stun level definitions
	// in this module.
	SetStun(lvl StunLevel)

	Blind() bool
	SetBlind(blind bool)

	Slow() bool
	SetSlow(slow bool)

	Confused() bool
	SetConfused(conf bool)

	Afraid() bool
	SetAfraid(fear bool)

	Paralyzed() bool
	SetParalyzed(para bool)

	Petrified() bool
	SetPetrified(para bool)

	Silenced() bool
	SetSilenced(sil bool)

	Cursed() bool
	SetCursed(c bool)

	// #blessed
	Blessed() bool
	SetBlessed(c bool)

	Corrosion() int
	SetCorrosion(level int)

	// Sight radius.
	Sight() int

	// Get this actor's current speed.
	// 1: Slow (0.5x normal)
	// 2: Normal
	// 3: Fast (1.5x normal)
	// 4: Very fast (2.x normal)
	Speed() int

	// Can this actor do stuff right now. This is false if they are paralyzed,
	// asleep, etc.
	CanAct() bool
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet struct {
	Trait
	stats  *stats
	skills *skills

	sight int
	speed int

	hp    int
	mp    int
	regen int

	stun      StunLevel
	corr      int
	blind     bool
	slow      bool
	afraid    bool
	confused  bool
	para      bool
	silence   bool
	cursed    bool
	blessed   bool
	petrified bool
}

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{
			stats: statlist{
				Str: 2,
				Agi: 5,
				Vit: 4,
				Mnd: 3,
			},
		},
		skills: &skills{},
		speed:  2,
		regen:  1,
		sight:  FOVRadius,
	}
	ps.hp = ps.MaxHP()
	ps.mp = ps.MaxMP()
	ps.initmods()
	return ps
}

// For testing.
func NewPlayerSheetFromSpec(pspec *PlayerSheet) *PlayerSheet {
	sheet := pspec.Copy()
	sheet.initmods()
	return sheet
}

func (p *PlayerSheet) Copy() *PlayerSheet {
	// Copy sheet.
	sheet := &PlayerSheet{}
	*sheet = *p

	// Copy stats.
	newstats := &stats{}
	if p.stats != nil {
		*newstats = *(p.stats)
	}
	sheet.stats = newstats

	newskills := &skills{}
	if p.skills != nil {
		*newskills = *(p.skills)
	}
	sheet.skills = newskills
	return sheet
}

func (p *PlayerSheet) Stat(stat StatName) int {
	return p.stats.stat(stat)
}

func (p *PlayerSheet) SetStat(stat StatName, amt int) {
	oldstat := p.Stat(stat)
	p.stats.set(stat, amt)
	newstat := p.Stat(stat)

	switch stat {
	case Agi:
		modAgiSkills(p, newstat-oldstat)
	case Vit:
		scaleHP(p, oldstat, newstat)
	case Mnd:
		modMndSkills(p, newstat-oldstat)
		scaleMP(p, oldstat, newstat)
	}

}

func (p *PlayerSheet) UnmodStat(stat StatName) int {
	return p.stats.unmodstat(stat)
}

func (p *PlayerSheet) StatMod(stat StatName) int {
	return p.stats.mod(stat)
}

func (p *PlayerSheet) ChangeStatMod(stat StatName, diff int) {
	old := p.Stat(stat)

	switch stat {
	case Agi:
		modAgiSkills(p, diff)
	case Vit:
		scaleHP(p, old, old+diff)
	case Mnd:
		modMndSkills(p, diff)
		scaleMP(p, old, old+diff)
	}
	p.stats.changemod(stat, diff)
}

func (p *PlayerSheet) Skill(skill SkillName) int {
	s := p.skills.skill(skill)
	if p.blind {
		s = blindpenalty(skill, s)
	}
	if !p.CanAct() {
		s = parapenalty(skill, s)
	}
	return s
}

func (p *PlayerSheet) SetSkill(skill SkillName, amt int) {
	p.skills.set(skill, amt)
}

func (p *PlayerSheet) UnmodSkill(skill SkillName) int {
	return p.skills.unmodskill(skill)
}

func (p *PlayerSheet) SkillMod(skill SkillName) int {
	return p.skills.mod(skill)
}

func (p *PlayerSheet) ChangeSkillMod(skill SkillName, diff int) {
	p.skills.changemod(skill, diff)
}

func (p *PlayerSheet) Speed() int {
	if p.slow {
		return slowpenalty(p.speed)
	}
	return p.speed
}

func (p *PlayerSheet) Dead() bool {
	return p.HP() <= 0
}

func (p *PlayerSheet) HP() int {
	return p.hp
}

func (p *PlayerSheet) setHP(hp int) {
	p.hp = hp
}

func (p *PlayerSheet) MaxHP() int {
	return vital(p.stats.stat(Vit))
}

func (p *PlayerSheet) MP() int {
	return p.mp
}

func (p *PlayerSheet) setMP(mp int) {
	p.mp = mp
}

func (p *PlayerSheet) MaxMP() int {
	return vital(p.stats.stat(Mnd))
}

func (p *PlayerSheet) Regen() int {
	if p.Petrified() {
		return 0
	}
	return p.regen
}

func (p *PlayerSheet) Sight() int {
	if p.blind {
		return 0
	}
	return p.sight
}

func (p *PlayerSheet) Hurt(dmg int) {
	hurt(p, dmg)
}

func (p *PlayerSheet) Heal(amt int) {
	heal(p, amt)
}

func (p *PlayerSheet) HurtMP(amt int) {
	hurtmp(p, amt)
}

func (p *PlayerSheet) HealMP(amt int) {
	healmp(p, amt)
}

func (p *PlayerSheet) Stun() StunLevel {
	return p.stun
}

func (p *PlayerSheet) SetStun(level StunLevel) {
	changestun(p, level)
	p.stun = level
}

func (p *PlayerSheet) SetBlind(b bool) {
	p.blind = b
}

func (p *PlayerSheet) Blind() bool {
	return p.blind
}

func (p *PlayerSheet) SetSlow(s bool) {
	p.slow = s
}

func (p *PlayerSheet) Slow() bool {
	return p.slow
}

func (p *PlayerSheet) SetAfraid(fear bool) {
	p.afraid = fear
}

func (p *PlayerSheet) Afraid() bool {
	return p.afraid
}

func (p *PlayerSheet) SetConfused(conf bool) {
	changeconf(p, conf)
	p.confused = conf
}

func (p *PlayerSheet) Confused() bool {
	return p.confused
}

func (p *PlayerSheet) SetParalyzed(para bool) {
	p.para = para
}

func (p *PlayerSheet) Paralyzed() bool {
	return p.para
}

func (p *PlayerSheet) SetSilenced(sil bool) {
	p.silence = sil
}

func (p *PlayerSheet) Silenced() bool {
	return p.silence
}

func (p *PlayerSheet) SetPetrified(pet bool) {
	p.petrified = pet
}

func (p *PlayerSheet) Petrified() bool {
	return p.petrified
}

func (p *PlayerSheet) SetCursed(c bool) {
	p.cursed = c
}

func (p *PlayerSheet) Cursed() bool {
	return p.cursed
}

func (p *PlayerSheet) SetBlessed(b bool) {
	p.blessed = b
}

func (p *PlayerSheet) Blessed() bool {
	return p.blessed
}

func (p *PlayerSheet) Corrosion() int {
	return p.corr
}

func (p *PlayerSheet) SetCorrosion(amt int) {
	p.corr = amt
}

func (p *PlayerSheet) CanAct() bool {
	return canact(p)
}

func (p *PlayerSheet) Attack() Attack {
	// We use skills.skill instead of Skill because Skill already applies the
	// blind penalty. We need to avoid it so we can apply the penalty to the
	// entire quantity.
	melee := p.obj.Equipper.Body().Melee() + p.skills.skill(Melee)
	if p.blind {
		melee = blindpenalty(Melee, melee)
	}

	weap := p.weapon()
	equip := weap.Equipment

	str := p.stats.stat(Str)
	bonusSides := math.Min(math.Abs(str), weap.Equipment.Weight) * math.Sgn(str)

	return Attack{
		Melee:   melee,
		Damroll: equip.Damroll.Add(0, bonusSides),
		CritDiv: equip.Weight + BaseCritDiv,
		Effects: equip.Effects,
		Verb:    "hits",
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
				Damroll: NewDice(1, math.Max(1, p.stats.stat(Str)+1)),
				Melee:   0,
				Weight:  0,
				Slot:    SlotHand,
			}),
		},
	})
}

func (p *PlayerSheet) Defense() Defense {
	body := p.obj.Equipper.Body()
	// We use skills.skill instead of Skill because Skill already applies the
	// blind penalty. We need to avoid it so we can apply the penalty to the
	// entire quantity.
	evasion := body.Evasion() + p.skills.skill(Evasion)
	if p.blind {
		evasion = blindpenalty(Evasion, evasion)
	}
	if !p.CanAct() {
		evasion = parapenalty(Evasion, evasion)
	}

	var dice []Dice
	effects := body.ArmorEffects()

	if p.Petrified() {
		dice = []Dice{petrifyProt}
		cr := NewEffects(map[Effect]int{ResistCrit: petrifyCritResist})
		effects = effects.Merge(cr)
	} else {
		dice = body.ProtDice()
	}

	corrdice := []Dice{}
	if corr := p.corr; corr > 0 {
		corrdice = append(corrdice, NewDice(corr, 4))
	}

	return Defense{
		Evasion:  evasion,
		ProtDice: dice,
		CorrDice: corrdice,
		Effects:  effects,
	}
}

func (p *PlayerSheet) initmods() {
	modAgiSkills(p, p.stats.stat(Agi))
	modMndSkills(p, p.stats.stat(Mnd))
}

// A spec for a monster attack.
type MonsterAttack struct {
	Attack
	// How relatively frequently should we use this attack?
	P int
}

func (m *MonsterAttack) Weight() int {
	return m.P
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	Trait

	stats  *stats
	skills *skills

	speed int
	sight int

	hp    int
	mp    int
	maxhp int
	maxmp int

	regen int

	stun      StunLevel
	blind     bool
	slow      bool
	confused  bool
	afraid    bool
	para      bool
	silence   bool
	cursed    bool
	blessed   bool
	petrified bool
	corr      int

	// The melee attacks this monster has. The elements in this slice should be
	// immutable -- do not change their fields, as these are shared across all
	// monsters that share the same spec.
	attacks []*MonsterAttack
	defense Defense
}

// Given a copy of a MonsterSheet literal, this will return a function that will bind
// the owner of the sheet to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewMonsterSheet(sheetspec *MonsterSheet) func(*Obj) Sheet {
	return func(o *Obj) Sheet {
		// Copy sheet.
		sheet := sheetspec.Copy()

		// Copy stats.
		sheet.obj = o
		sheet.hp = sheet.maxhp
		sheet.mp = sheet.maxmp
		if sheet.regen == 0 {
			sheet.regen = 1
		}
		if sheet.sight == 0 {
			sheet.sight = FOVRadius
		}
		return sheet
	}
}

// Make a deep copy of this sheet.
func (m *MonsterSheet) Copy() *MonsterSheet {
	// Copy sheet.
	sheet := &MonsterSheet{}
	*sheet = *m

	// Copy stats.
	newstats := &stats{}
	if m.stats != nil {
		*newstats = *(m.stats)
	}
	sheet.stats = newstats

	newskills := &skills{}
	if m.skills != nil {
		*newskills = *(m.skills)
	}
	sheet.skills = newskills
	return sheet
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

func (m *MonsterSheet) ChangeStatMod(stat StatName, diff int) {
	m.stats.changemod(stat, diff)
}

func (m *MonsterSheet) Skill(skill SkillName) int {
	if skill == Melee {
		panic("Monster Melee must be checked through individual attacks.")
	}
	if skill == Evasion {
		panic("Monster Evasion must be checked through Defense.")
	}
	s := m.skills.skill(skill)
	if !m.CanAct() {
		s = parapenalty(skill, s)
	}
	return s
}

func (m *MonsterSheet) SetSkill(skill SkillName, amt int) {
	m.skills.set(skill, amt)
}

func (m *MonsterSheet) UnmodSkill(skill SkillName) int {
	return m.skills.unmodskill(skill)
}

func (m *MonsterSheet) SkillMod(skill SkillName) int {
	return m.skills.mod(skill)
}

func (m *MonsterSheet) ChangeSkillMod(skill SkillName, diff int) {
	m.skills.changemod(skill, diff)
}

func (m *MonsterSheet) Speed() int {
	if m.Slow() {
		return slowpenalty(m.speed)
	}
	return m.speed
}

func (m *MonsterSheet) Dead() bool {
	return m.HP() <= 0
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
	if m.Petrified() {
		return 0
	}
	return m.regen
}

func (m *MonsterSheet) Sight() int {
	if m.blind {
		return 0
	}
	return m.sight
}

func (m *MonsterSheet) Hurt(dmg int) {
	hurt(m, dmg)
}

func (m *MonsterSheet) Heal(amt int) {
	heal(m, amt)
}

func (m *MonsterSheet) HurtMP(amt int) {
	hurtmp(m, amt)
}

func (m *MonsterSheet) HealMP(amt int) {
	healmp(m, amt)
}

func (m *MonsterSheet) Stun() StunLevel {
	return m.stun
}

func (m *MonsterSheet) SetStun(level StunLevel) {
	changestun(m, level)
	m.stun = level
}

func (m *MonsterSheet) SetAfraid(fear bool) {
	m.afraid = fear
}

func (m *MonsterSheet) Afraid() bool {
	return m.afraid
}

func (m *MonsterSheet) SetBlind(b bool) {
	m.blind = b
}

func (m *MonsterSheet) Blind() bool {
	return m.blind
}

func (m *MonsterSheet) SetSilenced(sil bool) {
	m.silence = sil
}

func (m *MonsterSheet) Silenced() bool {
	return m.silence
}

func (m *MonsterSheet) SetPetrified(p bool) {
	m.petrified = p
}

func (m *MonsterSheet) Petrified() bool {
	return m.petrified
}

func (m *MonsterSheet) SetSlow(s bool) {
	m.slow = s
}

func (m *MonsterSheet) Slow() bool {
	return m.slow
}

func (m *MonsterSheet) SetConfused(conf bool) {
	changeconf(m, conf)
	m.confused = conf
}

func (m *MonsterSheet) Confused() bool {
	return m.confused
}

func (m *MonsterSheet) SetParalyzed(b bool) {
	m.para = b
}

func (m *MonsterSheet) Paralyzed() bool {
	return m.para
}

func (p *MonsterSheet) SetCursed(c bool) {
	p.cursed = c
}

func (p *MonsterSheet) Cursed() bool {
	return p.cursed
}

func (p *MonsterSheet) SetBlessed(b bool) {
	p.blessed = b
}

func (p *MonsterSheet) Blessed() bool {
	return p.blessed
}

func (m *MonsterSheet) Corrosion() int {
	return m.corr
}

func (m *MonsterSheet) SetCorrosion(amt int) {
	m.corr = amt
}

func (m *MonsterSheet) CanAct() bool {
	return canact(m)
}

func (m *MonsterSheet) Attack() Attack {
	// Sigh.
	weighted := make([]Weighter, len(m.attacks))
	for i, at := range m.attacks {
		weighted[i] = at
	}

	pos, _ := WChoose(weighted)

	// Copy the attack.
	atk := m.attacks[pos].Attack
	atk.Melee = blindpenalty(Melee, atk.Melee)
	atk.CritDiv += BaseCritDiv
	return atk
}

func (m *MonsterSheet) Defense() Defense {
	def := m.defense

	if m.Petrified() {
		cr := NewEffects(map[Effect]int{ResistCrit: petrifyCritResist})
		def.Effects = def.Effects.Merge(cr)
		def.ProtDice = []Dice{petrifyProt}
	}

	def.CorrDice = []Dice{}
	if corr := m.corr; corr > 0 {
		def.CorrDice = append(def.CorrDice, NewDice(corr, 4))
	}

	return def
}

func modAgiSkills(s Sheet, diff int) {
	for sk := Melee; sk <= Stealth; sk++ {
		s.ChangeSkillMod(sk, diff)
	}
}

func modMndSkills(s Sheet, diff int) {
	for sk := Chi; sk < NumSkills; sk++ {
		s.ChangeSkillMod(sk, diff)
	}
}

func modAllSkills(s Sheet, diff int) {
	modAgiSkills(s, diff)
	modMndSkills(s, diff)
}

func modAllStats(s Sheet, diff int) {
	for stat := Str; stat < NumStats; stat++ {
		s.ChangeStatMod(stat, diff)
	}
}

// Maintains stats and modifications for a sheet.
type stats struct {
	stats statlist
	mods  statlist
}

// An array of statistics.
type statlist [NumStats]int

// Get the given stat combined with any active modifiers to it.
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

// Change a modifier for this stat by 'amt'.
func (s *stats) changemod(stat StatName, amt int) {
	s.mods[stat] += amt
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

// Maintains skills and modifications for a sheet. There are a number of ways
// that the stats and skills collections could have been combined into a single
// reusable data structure, but they were kept separate to make it very clear
// what is being used for what purpose, and especially so that indexing would
// be strongly typed, which otherwise would have been annoying to enforce.
type skills struct {
	skills skilllist
	mods   skilllist
}

// An array of skills
type skilllist [NumSkills]int

// Get the given skill combined with any active modifiers to it.
func (s *skills) skill(skill SkillName) int {
	return s.skills[skill] + s.mods[skill]
}

// Get the unmodified skill.
func (s *skills) unmodskill(skill SkillName) int {
	return s.skills[skill]
}

// Set the given skill to 'amt'.
func (s *skills) set(skill SkillName, amt int) {
	s.skills[skill] = amt
}

// Get the modifier (mod) for this skill.
func (s *skills) mod(skill SkillName) int {
	return s.mods[skill]
}

// Change a modifier for this skill by 'amt'.
func (s *skills) changemod(skill SkillName, amt int) {
	s.mods[skill] += amt
}

// Index into stats arrays.
type SkillName uint

// The actual stats our actors have.
const (
	// Called FIGHT in the game, but that is nasty to grep.
	Melee SkillName = iota
	// Called DODGE in the game, but also very verby.
	Evasion
	// SHOOT
	Shooting
	// SNEAK
	Stealth
	// CHI
	Chi
	// SENSE
	Sense
	// MAGIC
	Magic
	// SONG
	Song
	// Sentinel.
	NumSkills
)

// Details about an actor's melee attack, before the melee roll is applied --
// i.e. what melee bonus + damage should be done if no crits happen? What's the
// base divisor to use for calculating # of crits? What base verb should we use
// if it hits? (This latter may be altered by applybs() based on which effects
// actually end up working.)
type Attack struct {
	Melee   int
	Damroll Dice
	CritDiv int
	Effects Effects
	Verb    string
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
	CorrDice []Dice
	Effects  Effects
}

// Rolls protection dice - corrosion dice
func (def Defense) RollProt() int {
	sum := 0

	for _, d := range def.ProtDice {
		sum += d.Roll()
	}

	for _, d := range def.CorrDice {
		sum -= d.Roll()
	}

	return math.Max(sum, 0)
}

// Shared functions across sheets.
func heal(s Sheet, amt int) {
	s.setHP(math.Min(s.HP()+amt, s.MaxHP()))
}

func hurt(s Sheet, dmg int) {
	s.setHP(s.HP() - dmg)
	checkDeath(s)
}

func healmp(s Sheet, amt int) {
	s.setMP(math.Min(s.MP()+amt, s.MaxMP()))
}

func hurtmp(s Sheet, amt int) {
	s.setMP(math.Min(s.MP()-amt, 0))
	checkDeath(s)
}

func changestun(s Sheet, newstun StunLevel) {
	oldstun := s.Stun()
	if oldstun == newstun {
		return
	}

	modAllSkills(s, int(2*(oldstun-newstun)))

	msg := s.Obj().Spec.Name + " is "

	if oldstun > newstun {
		switch newstun {
		case Stunned:
			msg += "less stunned."
		case NotStunned:
			msg += "no longer stunned."
		}
	} else {
		switch newstun {
		case Stunned:
			msg += "stunned."
		case MoreStunned:
			msg += "very stunned."
		}
	}
	s.Obj().Game.Events.Message(msg)
}

func changeconf(s Sheet, newc bool) {
	oldc := s.Confused()
	if oldc == newc {
		return
	} else if oldc == true {
		modMndSkills(s, 5)
	} else {
		modMndSkills(s, dumpsterEvasion)
	}
}

func checkDeath(s Sheet) {
	if !s.Dead() {
		return
	}

	obj := s.Obj()
	game := obj.Game

	if obj.Dropper != nil {
		obj.Dropper.DropItems()
	}

	game.Events.Message(fmt.Sprintf("%s fell.", obj.Spec.Name))
	game.Kill(obj)
}

func canact(s Sheet) bool {
	return !((s.Paralyzed() || s.Petrified()) && !s.Dead())
}

func blindpenalty(skill SkillName, score int) int {
	if skill == Melee || skill == Evasion || skill == Shooting {
		return score / 2
	}
	return score
}

func parapenalty(skill SkillName, score int) int {
	if skill == Evasion {
		return dumpsterEvasion
	}
	return score
}

func slowpenalty(spd int) int {
	return math.Max(spd-1, 1)
}

func vital(stat int) int {
	return 10 * (1 + math.Max(stat, 1))
}

func scaleHP(sheet Sheet, oldv, newv int) {
	sheet.setHP(sheet.HP() * vital(newv) / vital(oldv))
}

func scaleMP(sheet Sheet, oldv, newv int) {
	sheet.setMP(sheet.MP() * vital(newv) / vital(oldv))
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
)

// The base divisor to use for crits.
const (
	BaseCritDiv       = 7
	dumpsterEvasion   = -5
	petrifyCritResist = 100
)

var petrifyProt = NewDice(8, 4)
