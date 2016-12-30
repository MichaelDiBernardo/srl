package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"math/rand"
)

// Permissible genera of actors.
const (
	GenPlayer  = "player"
	GenMonster = "monster"
)

// The base divisor to use for crits.
const baseCritDiv = 7

// A thing that can move given a specific direction.
type Mover interface {
	Objgetter
	Move(dir math.Point) bool
}

// A universally-applicable mover for actors.
type ActorMover struct {
	Trait
}

// Constructor for actor movers.
func NewActorMover(obj *Obj) Mover {
	return &ActorMover{Trait: Trait{obj: obj}}
}

// Try to move the actor. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	endtile := obj.Level.At(endpos)
	if other := endtile.Actor; other != nil {
		if opposing := obj.IsPlayer() != other.IsPlayer(); opposing {
			p.obj.Fighter.Hit(other.Fighter)
		}
		return false
	}

	moved := obj.Level.Place(obj, endpos)
	if moved {
		if items := endtile.Items; !items.Empty() && obj.IsPlayer() {
			var msg string
			topname, n := items.Top().Spec.Name, items.Len()
			if n == 1 {
				msg = fmt.Sprintf("%v sees %v here.", obj.Spec.Name, topname)
			} else {
				msg = fmt.Sprintf("%v sees %v and %d other items here.", obj.Spec.Name, topname, n-1)
			}
			obj.Game.Events.Message(msg)
		}
	}
	return moved
}

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Act(l *Level) bool
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	Trait
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{Trait: Trait{obj: obj}}
}

// Move in any of the 8 directions with uniform chance. Does not take walls
// etc. in account so this will happily try to bump into things.
func (ai *RandomAI) Act(l *Level) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	dir := math.Pt(x, y)
	if dir == math.Origin {
		return ai.Act(l)
	}
	return ai.obj.Mover.Move(dir)
}

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

func (p *PlayerSheet) Attack() Attack {
	melee := p.obj.Equipper.Body().Melee() + p.Melee()

	weap := p.weapon()
	str := p.Str()

	bonusSides := math.Min(math.Abs(str), weap.Equip.Weight) * math.Sgn(str)
	return Attack{
		Melee:   melee,
		Damroll: weap.Equip.Damroll.Add(0, bonusSides),
		CritDiv: weap.Equip.Weight + baseCritDiv,
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
		Genus:   GenEquip,
		Species: SpecFist,
		Name:    "FIST",
		Traits: &Traits{
			Equip: NewEquip(Equip{
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

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
}

func hit(attacker Fighter, defender Fighter) {
	atk, def := attacker.Obj().Sheet.Attack(), defender.Obj().Sheet.Defense()
	residual := DieRoll(1, 20) + atk.Melee - DieRoll(1, 20) + def.Evasion
	aname, dname := attacker.Obj().Spec.Name, defender.Obj().Spec.Name

	if residual > 0 {
		crits := residual / atk.CritDiv
		dmg := math.Max(0, atk.RollDamage(crits)-def.RollProt())

		critstr := ""
		if crits > 0 {
			critstr = fmt.Sprintf(" %dx critical!", crits)
		}
		msg := fmt.Sprintf("%v hit %v (%d).%s", aname, dname, dmg, critstr)
		attacker.Obj().Game.Events.Message(msg)

		defender.Obj().Sheet.Hurt(dmg)
	} else {
		msg := fmt.Sprintf("%v missed %v.", aname, dname)
		attacker.Obj().Game.Events.Message(msg)
	}
}

// Melee combat.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *ActorFighter) Hit(other Fighter) {
	hit(f, other)
}

// A thing that that can hold items in inventory. (A "pack".)
type Packer interface {
	Objgetter
	// Tries to pickup something at current square. If there are many things,
	// will invoke stack menu.
	TryPickup()
	// Pickup the item on the floor stack at given index.
	Pickup(index int)
	// Tries to drop something at current square.
	TryDrop()
	// Drop the item at index in inventory to the floor stack.
	Drop(index int)
	// Get this Packer's inventory.
	Inventory() *Inventory
}

// An attacker that works for all actors.
type ActorPacker struct {
	Trait
	inventory *Inventory
}

func NewActorPacker(obj *Obj) Packer {
	return &ActorPacker{
		Trait:     Trait{obj: obj},
		inventory: NewInventory(),
	}
}

func (a *ActorPacker) Inventory() *Inventory {
	return a.inventory
}

func (a *ActorPacker) TryPickup() {
	ground := a.obj.Tile.Items
	if ground.Empty() {
		a.obj.Game.Events.Message("Nothing there.")
	} else if ground.Len() == 1 {
		a.moveFromGround(0)
	} else {
		a.obj.Game.SwitchMode(ModePickup)
	}
}

func (a *ActorPacker) Pickup(index int) {
	a.obj.Game.SwitchMode(ModeHud)
	a.moveFromGround(index)
}

func (a *ActorPacker) TryDrop() {
	if a.inventory.Empty() {
		a.obj.Game.Events.Message("Nothing to drop.")
		return
	}

	ground := a.obj.Tile.Items
	if ground.Full() {
		a.obj.Game.Events.Message("Can't drop here.")
	} else {
		a.obj.Game.SwitchMode(ModeDrop)
	}
}

func (a *ActorPacker) Drop(index int) {
	a.obj.Game.SwitchMode(ModeHud)
	item := a.inventory.Take(index)

	// Bounds-check the index the player requested.
	if item == nil {
		return
	}
	a.obj.Tile.Items.Add(item)
	a.obj.Game.Events.Message(fmt.Sprintf("%v dropped %v.", a.obj.Spec.Name, item.Spec.Name))
}

func (a *ActorPacker) moveFromGround(index int) {
	// Bounds-check the index the player requested.
	item := a.obj.Tile.Items.At(index)
	if item == nil {
		return
	}

	if a.inventory.Full() {
		a.obj.Game.Events.Message(fmt.Sprintf("%v has no room for %v.", a.obj.Spec.Name, item.Spec.Name))
	} else {
		item := a.obj.Tile.Items.Take(index)
		a.inventory.Add(item)
		a.obj.Game.Events.Message(fmt.Sprintf("%v got %v.", a.obj.Spec.Name, item.Spec.Name))
	}
}

type Equipper interface {
	Objgetter
	// Bring up the equipper screen if anything in inventory can be equipped.
	TryEquip()
	// Bring up the remover screen if anything on body can be removed.
	TryRemove()
	// Equip the item at index 'index' in inventory.
	Equip(index int)
	// Remove the item equipped in the given slot.
	Remove(slot Slot)
	// Get the underlying entity's Body.
	Body() *Body
}

type ActorEquipper struct {
	Trait
	body *Body
}

func NewActorEquipper(obj *Obj) Equipper {
	return &ActorEquipper{
		Trait: Trait{obj: obj},
		body:  NewBody(),
	}
}

func (a *ActorEquipper) TryEquip() {
	if !a.obj.Packer.Inventory().HasEquips() {
		a.obj.Game.Events.Message("Nothing to wield/wear.")
	} else {
		a.obj.Game.SwitchMode(ModeEquip)
	}
}

func (a *ActorEquipper) TryRemove() {
	if a.body.Naked() {
		a.obj.Game.Events.Message("Not wearing anything.")
	} else if a.obj.Packer.Inventory().Full() && a.obj.Tile.Items.Full() {
		a.obj.Game.Events.Message("Can't remove; pack and ground are full.")
	} else {
		a.obj.Game.SwitchMode(ModeRemove)
	}
}

func (a *ActorEquipper) Equip(index int) {
	a.obj.Game.SwitchMode(ModeHud)

	equip := a.obj.Packer.Inventory().Take(index)

	// Bounds-check the index the player requested.
	if equip == nil {
		return
	}

	if equip.Spec.Genus != GenEquip {
		a.obj.Game.Events.Message(fmt.Sprintf("Cannot equip %v.", equip.Spec.Name))
		return
	}

	if swapped := a.body.Wear(equip); swapped != nil {
		a.obj.Packer.Inventory().Add(swapped)
	}
}

func (a *ActorEquipper) Remove(slot Slot) {
	a.obj.Game.SwitchMode(ModeHud)

	removed := a.body.Remove(slot)
	if removed == nil {
		return
	}

	if added := a.obj.Packer.Inventory().Add(removed); added {
		return
	}

	// No room for unequipped item in inventory; drop it.
	a.obj.Tile.Items.Add(removed)
	a.obj.Game.Events.Message(fmt.Sprintf("No room in pack! Dropped %v.", removed.Spec.Name))
}

func (a *ActorEquipper) Body() *Body {
	return a.body
}
