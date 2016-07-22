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
	Actors []*Obj
}

// Create a level generated by the given generator function.
func NewLevel(width, height int, gen func(*Level) *Level) *Level {
	newmap := Map{}
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		newmap = append(newmap, row)
	}
	level := &Level{Map: newmap, Actors: make([]*Obj, 10)}
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
