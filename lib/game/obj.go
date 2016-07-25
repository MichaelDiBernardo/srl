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

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	Type    ObjType
	Subtype ObjSubtype
	Tile    *Tile
	Level   *Level
	Events  *EventQueue
	// Actor traits
	Mover Mover
	AI    AI
}

// A specification object for newObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover func(*Obj) Mover
	AI    func(*Obj) AI
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
	return t
}

// Create a new game object with the given traits. This shouldn't be used
// directly; you should instead use a *Game as a factory for any game objects
// that need creating.
func newObj(otype ObjType, ostype ObjSubtype, traits *Traits, name string, eq *EventQueue) *Obj {
	// Create.
	newobj := &Obj{Type: otype, Subtype: ostype, Events: eq}

	// Assign traits.
	traits = traits.defaults()
	newobj.Mover = traits.Mover(newobj)
	newobj.AI = traits.AI(newobj)

	return newobj
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}
