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

// A dummy mover used in cases where a thing can't move.
type nullMover struct {
	Trait
}

// Do nothing and return false.
func (_ *nullMover) Move(dir math.Point) bool {
	return false
}

// Constructor for null movers.
func NewNullMover(obj *Obj) Mover {
	return &nullMover{Trait: Trait{obj: obj}}
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

// A dummy AI used when a thing doesn't need a computer to think for it.
type nullAI struct {
	Trait
}

// Do nothing and return false.
func (_ *nullAI) Act(l *Level) bool {
	return false
}

// Constructor for null movers.
func NewNullAI(obj *Obj) AI {
	return &nullAI{Trait: Trait{obj: obj}}
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

// Accessors for an actor's stats.
type Stats interface {
	Objgetter
	Str() int
	Agi() int
	Vit() int
	Mnd() int
}

// Single implementation of this for now; will probably have separate
// implementations for monsters and player when things get more complicated.
type stats struct {
	Trait
	str int
	agi int
	vit int
	mnd int
}

// Given a copy of a stats literal, this will return a function that will bind
// the owner of the stats to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewActorStats(stats stats) func(*Obj) Stats {
	return func(o *Obj) Stats {
		stats.obj = o
		return &stats
	}
}

func (s *stats) Str() int {
	return s.str
}

func (s *stats) Agi() int {
	return s.agi
}

func (s *stats) Vit() int {
	return s.vit
}

func (s *stats) Mnd() int {
	return s.mnd
}

func NewNullStats(obj *Obj) Stats {
	return &stats{Trait: Trait{obj: obj}}
}

// A 'character sheet' for an actor. This is where all attributes derived from
// stats + equipment are stored.
type Sheet interface {
	Objgetter
	Melee() int
	Evasion() int
	HP() int
	MaxHP() int
	MP() int
	MaxMP() int

	Hurt(dmg int)
}

type sheet struct {
	Trait
	hp int
	mp int
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet sheet

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{Trait: Trait{obj: obj}}
	ps.hp = ps.MaxHP()
	ps.mp = ps.MaxMP()
	return ps
}

func (p *PlayerSheet) Melee() int {
	return p.obj.Stats.Agi()
}

func (p *PlayerSheet) Evasion() int {
	return p.obj.Stats.Agi()
}

func (p *PlayerSheet) HP() int {
	return p.hp
}

func (p *PlayerSheet) MP() int {
	return p.mp
}

func (p *PlayerSheet) MaxHP() int {
	return 10 * (1 + p.obj.Stats.Vit())
}

func (p *PlayerSheet) MaxMP() int {
	return 10 * (1 + p.obj.Stats.Mnd())
}

func (p *PlayerSheet) Hurt(dmg int) {
	p.hp -= dmg
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	Trait
	sheet
	melee   int
	evasion int
	maxhp   int
	maxmp   int
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

func NewNullSheet(obj *Obj) Sheet {
	return &MonsterSheet{Trait: Trait{obj: obj}}
}

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
	MeleeRoll() int
	EvasionRoll() int
	DamRoll() int
	ProtRoll() int
}

// An attacker that works for all actors.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{Trait: Trait{obj: obj}}
}

func (a *ActorFighter) Hit(other Fighter) {
	mroll, eroll := a.MeleeRoll(), other.EvasionRoll()
	if mroll > eroll {
		dmg := a.DamRoll() - other.ProtRoll()
		other.Obj().Sheet.Hurt(dmg)
		msg := fmt.Sprintf("%v hit %v (%d).", a.obj.Spec.Name, other.Obj().Spec.Name, dmg)
		a.obj.Game.Events.Message(msg)
	} else {
		msg := fmt.Sprintf("%v missed %v.", a.obj.Spec.Name, other.Obj().Spec.Name)
		a.obj.Game.Events.Message(msg)
	}
}

func (a *ActorFighter) MeleeRoll() int {
	return DieRoll(1, 20) + a.obj.Sheet.Melee()
}

func (a *ActorFighter) EvasionRoll() int {
	return DieRoll(1, 20) + a.obj.Sheet.Evasion()
}

func (a *ActorFighter) DamRoll() int {
	return DieRoll(1, a.obj.Stats.Str())
}

func (a *ActorFighter) ProtRoll() int {
	return 0
}

type NullFighter struct {
	Trait
}

func NewNullFighter(obj *Obj) Fighter {
	return &NullFighter{Trait: Trait{obj: obj}}
}

func (n *NullFighter) Hit(other Fighter) {
}

func (n *NullFighter) MeleeRoll() int {
	return 0
}

func (n *NullFighter) EvasionRoll() int {
	return 0
}

func (n *NullFighter) DamRoll() int {
	return 0
}

func (n *NullFighter) ProtRoll() int {
	return 0
}

// A thing that that can hold items in inventory. (A "pack".)
type Packer interface {
	Objgetter
	// Tries to pickup something at current square. If there are many things,
	// will invoke stack menu.
	TryPickup()
	// Pickup the item on the floor stack at given index. Returns false if
	// there was no room in player inventory to do this.
	Pickup(index int)
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
		a.obj.Game.Events.Message(fmt.Sprintf("Nothing there."))
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

func (a *ActorPacker) moveFromGround(index int) {
	if a.inventory.Full() {
		item := a.obj.Tile.Items.At(index)
		a.obj.Game.Events.Message(fmt.Sprintf("%v has no room for %v.", a.obj.Spec.Name, item.Spec.Name))
	} else {
		item := a.obj.Tile.Items.Take(index)
		a.inventory.Add(item)
		a.obj.Game.Events.Message(fmt.Sprintf("%v got %v.", a.obj.Spec.Name, item.Spec.Name))
	}
}

type NullPacker struct {
	Trait
}

func NewNullPacker(obj *Obj) Packer {
	return &NullPacker{Trait: Trait{obj: obj}}
}

func (a *NullPacker) TryPickup() {
}

func (a *NullPacker) Inventory() *Inventory {
	return nil
}

func (a *NullPacker) Pickup(index int) {
}
