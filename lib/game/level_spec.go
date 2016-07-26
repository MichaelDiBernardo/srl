package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

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

func TestLevel(l *Level) *Level {
	l = SquareLevel(l)
	l.Place(l.fac.NewObj(MonOrc), math.Pt(10, 10))
	l.Place(l.fac.NewObj(MonOrc), math.Pt(20, 20))
	return l
}
