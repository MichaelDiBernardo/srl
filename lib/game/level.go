package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type Tile struct {
	Feature *Feature
	Actor   *Obj
	Pos     math.Point
}

type Map [][]*Tile

type Level struct {
	Map    Map
	Bounds math.Rectangle
	fac    ObjFactory
	actors []*Obj
}

// Create a level that uses the given factory to create game objects, and which
// generated by the given generator function.
func NewLevel(width, height int, fac ObjFactory, gen func(*Level) *Level) *Level {
	newmap := Map{}
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		newmap = append(newmap, row)
	}
	level := &Level{
		Map:    newmap,
		Bounds: math.Rect(math.Origin, math.Pt(width, height)),
		fac:    fac,
		actors: make([]*Obj, 0),
	}
	return gen(level)
}

func (l *Level) At(p math.Point) *Tile {
	return l.Map[p.Y][p.X]
}

func (l *Level) HasPoint(p math.Point) bool {
	return l.Bounds.HasPoint(p)
}

// Place `o` on the tile at `p`. Returns false if this is impossible (e.g.
// trying to put something on a solid square.)
// This will remove `o` from any tile on any map it was previously on.
func (l *Level) Place(o *Obj, p math.Point) bool {
	tile := l.At(p)

	if tile.Feature.Solid || tile.Actor != nil {
		return false
	}

	// If this actor has been placed before, we need to clear the tile they
	// were on previously. If they haven't, we need to add them to the actor
	// list so we know who they are.
	if o.Tile != nil {
		o.Tile.Actor = nil
	} else {
		l.actors = append(l.actors, o)
	}

	o.Level = l
	o.Tile = tile

	tile.Actor = o

	return true
}

func (l *Level) Evolve() {
	for _, actor := range l.actors {
		actor.AI.Act(l)
	}
}
