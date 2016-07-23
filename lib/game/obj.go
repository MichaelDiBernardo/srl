package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// e.g. "Actor", "Item"
type ObjType string

// e.g. "Orc", "HealPotion"
type ObjSubtype string

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	Type    ObjType
	Subtype ObjSubtype
	Tile    *Tile
	Level   *Level
	// Actor traits
	Mover Mover
	AI    AI
}

// A specification object for NewObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover func(*Obj) Mover
	AI    func(*Obj) AI
}

// Create a new game object with the given traits. This object should be placed
// on a map with Place before use.
func NewObj(traits *Traits) *Obj {
	newobj := &Obj{}
	if traits.Mover != nil {
		newobj.Mover = traits.Mover(newobj)
	}
	if traits.AI != nil {
		newobj.AI = traits.AI(newobj)
	}
	return newobj
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}
