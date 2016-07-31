package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// e.g. "Actor", "Item"
type ObjType string

// e.g. "Orc", "HealPotion"
type ObjSubtype string

const (
	OTActor = "Actor"
)

// A specification for a type of game object.
type Spec struct {
	Type    ObjType
	Subtype ObjSubtype
	Name    string
	Traits  *Traits
}

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	Spec   *Spec
	Tile   *Tile
	Level  *Level
	Events *EventQueue
	// Actor traits
	Mover Mover
	AI    AI
	Stats Stats
	Sheet Sheet
}

// A specification object for newObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover func(*Obj) Mover
	AI    func(*Obj) AI
	Stats func(*Obj) Stats
	Sheet func(*Obj) Sheet
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
	return t
}

// Create a new game object with the given traits. This shouldn't be used
// directly; you should instead use a *Game as a factory for any game objects
// that need creating.
func newObj(spec *Spec, eq *EventQueue) *Obj {
	// Create.
	newobj := &Obj{Spec: spec, Events: eq}

	// Assign traits.
	traits := spec.Traits.defaults()
	newobj.Mover = traits.Mover(newobj)
	newobj.AI = traits.AI(newobj)
	newobj.Stats = traits.Stats(newobj)
	newobj.Sheet = traits.Sheet(newobj)

	return newobj
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}
