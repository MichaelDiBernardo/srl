package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// e.g. "Actor", "Item"
type Family string

// e.g. "Player", "Monster", "Potion"
type Genus string

// e.g. "Orc", "HealPotion"
type Species string

// Permissible families of objects.
const (
	FamActor = "actor"
	FamItem  = "item"
)

// Trait basics; the ability to backreference the attached object.
type Objgetter interface {
	Obj() *Obj
}

type Trait struct {
	obj *Obj
}

func (t *Trait) Obj() *Obj {
	return t.obj
}

// A specification for a type of game object.
type Spec struct {
	Family  Family
	Genus   Genus
	Species Species
	Name    string
	Traits  *Traits
}

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	Spec *Spec

	Tile  *Tile
	Level *Level
	Game  *Game

	Mover   Mover
	AI      AI
	Stats   Stats
	Sheet   Sheet
	Fighter Fighter
	Packer  Packer
}

func (o *Obj) String() string {
	if o.Spec != nil {
		return "Obj: " + o.Spec.Name
	}
	return "Obj: ???"
}

// A specification object for newObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover   func(*Obj) Mover
	AI      func(*Obj) AI
	Stats   func(*Obj) Stats
	Sheet   func(*Obj) Sheet
	Fighter func(*Obj) Fighter
	Packer  func(*Obj) Packer
}

// Takes a partially-specified traits obj and fills in the nil ones with
// nullobj versions.
func (t *Traits) defaults() *Traits {
	if t.Mover == nil {
		t.Mover = NewNullMover
	}
	if t.AI == nil {
		t.AI = NewNullAI
	}
	if t.Stats == nil {
		t.Stats = NewNullStats
	}
	if t.Sheet == nil {
		t.Sheet = NewNullSheet
	}
	if t.Fighter == nil {
		t.Fighter = NewNullFighter
	}
	if t.Packer == nil {
		t.Packer = NewNullPacker
	}
	return t
}

// Create a new game object from the given spec. This shouldn't be used
// directly; you should instead use a *Game as a factory for any game objects
// that need creating. This will not initialize the fields on the obj that have
// nothing to do with specs or traits (e.g. game, eventqueue, tile etc.)
func newObj(spec *Spec) *Obj {
	// Create.
	newobj := &Obj{Spec: spec}

	// Assign traits.
	traits := spec.Traits.defaults()
	newobj.Mover = traits.Mover(newobj)
	newobj.AI = traits.AI(newobj)
	newobj.Stats = traits.Stats(newobj)
	newobj.Sheet = traits.Sheet(newobj)
	newobj.Fighter = traits.Fighter(newobj)
	newobj.Packer = traits.Packer(newobj)

	return newobj
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}

// Does this object represent the player?
func (o *Obj) isPlayer() bool {
	return o.Spec.Family == FamActor && o.Spec.Genus == GenPlayer
}
