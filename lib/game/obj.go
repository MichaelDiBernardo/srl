package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	Tile  *Tile
	Map   Map
	Mover Mover
}

// A specification object for NewObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover func(*Obj) Mover
}

// Create a new game object with the given traits. This object should be placed
// on a map with Place before use.
func NewObj(traits Traits) *Obj {
	newobj := &Obj{}
    if traits.Mover != nil {
        newobj.Mover = traits.Mover(newobj)
    }
	return newobj
}

// Place `o` on the tile at `p`. Returns false if this is impossible (e.g.
// trying to put something on a solid square.)
// This will remove `o` from any tile on any map it was previously on.
func (o *Obj) Place(m Map, p math.Point) bool {
	tile := m.At(p)

	if tile.Feature.Solid {
		return false
	}

	if o.Tile != nil {
		o.Tile.Actor = nil
	}
	o.Map = m
	o.Tile = tile
	tile.Actor = o
	return true
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}
