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
		if opposing := obj.isPlayer() != other.isPlayer(); opposing {
			p.obj.Fighter.Hit(other.Fighter)
		}
		return false
	}

	moved := obj.Level.Place(obj, endpos)
	if moved {
		if items := endtile.Items; !items.Empty() && obj.isPlayer() {
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
		str:   2,
		agi:   2,
		vit:   2,
		mnd:   2,
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

	protroll Dice
	damroll  Dice
}

// Given a copy of a MonsterSheet literal, this will return a function that will bind
// the owner of the sheet to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewMonsterSheet(sheet MonsterSheet) func(*Obj) Sheet {
	return func(o *Obj) Sheet {
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
}

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
	MeleeRoll() int
	EvasionRoll() int
	Damroll() int
	Protroll() int
}

func hit(attacker Fighter, defender Fighter) {
	mroll, eroll := attacker.MeleeRoll(), defender.EvasionRoll()
	if mroll > eroll {
		dmg := attacker.Damroll() - defender.Protroll()
		defender.Obj().Sheet.Hurt(dmg)
		msg := fmt.Sprintf("%v hit %v (%d).", attacker.Obj().Spec.Name, defender.Obj().Spec.Name, dmg)
		attacker.Obj().Game.Events.Message(msg)
	} else {
		msg := fmt.Sprintf("%v missed %v.", attacker.Obj().Spec.Name, defender.Obj().Spec.Name)
		attacker.Obj().Game.Events.Message(msg)
	}
}

// Player melee combat.
type PlayerFighter struct {
	Trait
}

func NewPlayerFighter(obj *Obj) Fighter {
	return &PlayerFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *PlayerFighter) Hit(other Fighter) {
	hit(f, other)
}

func (f *PlayerFighter) MeleeRoll() int {
	bonus := f.obj.Equipper.Body().Melee() + f.obj.Sheet.Melee()
	return DieRoll(1, 20) + bonus
}

func (f *PlayerFighter) EvasionRoll() int {
	bonus := f.obj.Equipper.Body().Evasion() + f.obj.Sheet.Evasion()
	return DieRoll(1, 20) + bonus
}

func (f *PlayerFighter) Damroll() int {
	weap := f.weapon()
	str := f.obj.Sheet.Str()
	bonus := math.Min(math.Abs(str), weap.Equip.Weight) * math.Sgn(str)
	return weap.Equip.Damroll.Add(0, bonus).Roll()
}

func (f *PlayerFighter) Protroll() int {
	dice := f.obj.Equipper.Body().ProtDice()
	sum := 0

	for i := 0; i < len(dice); i++ {
		sum += dice[i].Roll()
	}

	return sum
}

func (f *PlayerFighter) weapon() *Obj {
	weap := f.obj.Equipper.Body().Weapon()
	if weap != nil {
		return weap
	}
	return f.fist()
}

func (f *PlayerFighter) fist() *Obj {
	return f.obj.Game.NewObj(&Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: SpecFist,
		Name:    "FIST",
		Traits: &Traits{
			Equip: NewEquip(Equip{
				Damroll: NewDice(1, f.obj.Sheet.Str()+1),
				Melee:   0,
				Weight:  0,
				Slot:    SlotHand,
			}),
		},
	})
}

// Monster melee combat.
type MonsterFighter struct {
	Trait
}

func NewMonsterFighter(obj *Obj) Fighter {
	return &MonsterFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *MonsterFighter) Hit(other Fighter) {
	mroll, eroll := f.MeleeRoll(), other.EvasionRoll()
	if mroll > eroll {
		dmg := f.Damroll() - other.Protroll()
		other.Obj().Sheet.Hurt(dmg)
		msg := fmt.Sprintf("%v hit %v (%d).", f.obj.Spec.Name, other.Obj().Spec.Name, dmg)
		f.obj.Game.Events.Message(msg)
	} else {
		msg := fmt.Sprintf("%v missed %v.", f.obj.Spec.Name, other.Obj().Spec.Name)
		f.obj.Game.Events.Message(msg)
	}
}

func (f *MonsterFighter) MeleeRoll() int {
	return DieRoll(1, 20) + f.obj.Sheet.Melee()
}

func (f *MonsterFighter) EvasionRoll() int {
	return DieRoll(1, 20) + f.obj.Sheet.Evasion()
}

func (f *MonsterFighter) Damroll() int {
	return f.obj.Sheet.(*MonsterSheet).damroll.Roll()
}

func (f *MonsterFighter) Protroll() int {
	return f.obj.Sheet.(*MonsterSheet).protroll.Roll()
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
