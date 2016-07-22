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

// Create a map of the given height and width that will have its tiles
// generated by the given generator function.
func NewMap(width, height int, gen func(Map) Map) Map {
	newmap := Map{}
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		newmap = append(newmap, row)
	}
    return gen(newmap)
}

func (m Map) Width() int {
	return len(m[0])
}

func (m Map) Height() int {
	return len(m)
}

func (m Map) At(p math.Point) *Tile {
	return m[p.Y][p.X]
}

func (m Map) HasPoint(p math.Point) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < m.Width() && p.Y < m.Height()
}

// Generators.
func SquareMap(m Map) Map {
    height, width := m.Height(), m.Width()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			feature := FeatFloor
			if x == 0 || y == 0 || y == height-1 || x == width-1 {
				feature = FeatWall
			}
			m[y][x] = &Tile{Pos: math.Pt(x, y), Feature: feature}
		}
	}
    return m
}

func IdentMap(m Map) Map {
    return m
}

