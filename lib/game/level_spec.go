package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"math/rand"
)

// Generators.
func SquareLevel(l *Level) *Level {
	height, width, m := l.Bounds.Height(), l.Bounds.Width(), l.Map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			feature := FeatFloor
			if x == 0 || y == 0 || y == height-1 || x == width-1 {
				feature = FeatWall
			}
			m[y][x].Feature = feature
		}
	}
	return l
}

func IdentLevel(l *Level) *Level {
	height, width, m := l.Bounds.Height(), l.Bounds.Width(), l.Map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m[y][x].Feature = FeatFloor
		}
	}
	return l
}

func TestLevel(l *Level) *Level {
	l = SquareLevel(l)
	for i := 0; i < 10; i++ {
		mon := l.fac.NewObj(MonOrc)
		for {
			pt := math.Pt(rand.Intn(l.Bounds.Width()), rand.Intn(l.Bounds.Height()))
			if l.Place(mon, pt) {
				break
			}
		}
	}
	for i := 0; i < 10; i++ {
		item := l.fac.NewObj(WeapSword)
		for {
			pt := math.Pt(rand.Intn(l.Bounds.Width()), rand.Intn(l.Bounds.Height()))
			if l.Place(item, pt) {
				if rand.Intn(2) == 1 {
					extras := rand.Intn(4) + 1
					for j := 0; j < extras; j++ {
						l.Place(l.fac.NewObj(WeapSword), pt)
					}
				}
				break
			}
		}
	}
	return l
}
