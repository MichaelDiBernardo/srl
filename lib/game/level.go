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
	actors []*Obj
}

// Create a level generated by the given generator function.
func NewLevel(width, height int, gen func(*Level) *Level) *Level {
	newmap := Map{}
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		newmap = append(newmap, row)
	}
	level := &Level{Map: newmap, actors: make([]*Obj, 0)}
	return gen(level)
}

func (l *Level) Width() int {
	return len(l.Map[0])
}

func (l *Level) Height() int {
	return len(l.Map)
}

func (l *Level) At(p math.Point) *Tile {
	return l.Map[p.Y][p.X]
}

func (l *Level) HasPoint(p math.Point) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < l.Width() && p.Y < l.Height()
}

// Place `o` on the tile at `p`. Returns false if this is impossible (e.g.
// trying to put something on a solid square.)
// This will remove `o` from any tile on any map it was previously on.
func (l *Level) Place(o *Obj, p math.Point) bool {
	tile := l.At(p)

	if tile.Feature.Solid {
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
		if actor.AI != nil {
			actor.AI.Act(l)
		}
	}
}

func TestLevel(l *Level) *Level {
	l = SquareLevel(l)
	l.Place(NewActor(MonOrc), math.Pt(10, 10))
	l.Place(NewActor(MonOrc), math.Pt(20, 20))
	return l
}

// Generators.
func SquareLevel(l *Level) *Level {
	height, width, m := l.Height(), l.Width(), l.Map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			feature := FeatFloor
			if x == 0 || y == 0 || y == height-1 || x == width-1 {
				feature = FeatWall
			}
			m[y][x] = &Tile{Pos: math.Pt(x, y), Feature: feature}
		}
	}
	return l
}

func IdentLevel(l *Level) *Level {
	height, width, m := l.Height(), l.Width(), l.Map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m[y][x] = &Tile{Pos: math.Pt(x, y), Feature: FeatFloor}
		}
	}
	return l
}
