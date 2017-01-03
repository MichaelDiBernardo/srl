package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
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
	l.Place(l.game.Player, math.Pt(1, 1))
	return l
}

func IdentLevel(l *Level) *Level {
	height, width, m := l.Bounds.Height(), l.Bounds.Width(), l.Map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m[y][x].Feature = FeatFloor
		}
	}
	l.Place(l.game.Player, math.Pt(1, 1))
	return l
}

func TestLevel(l *Level) *Level {
	l = SquareLevel(l)
	for i := 0; i < 10; i++ {
		mon := l.game.NewObj(MonOrc)
		for {
			pt := math.Pt(rand.Intn(l.Bounds.Width()), rand.Intn(l.Bounds.Height()))
			if l.Place(mon, pt) {
				break
			}
		}
	}
	for i := 0; i < 20; i++ {
		var pt math.Point
		for {
			pt = math.Pt(rand.Intn(l.Bounds.Width()), rand.Intn(l.Bounds.Height()))
			if !l.At(pt).Feature.Solid {
				break
			}
		}

		stacksize := rand.Intn(3) + 1
		for j := 0; j < stacksize; j++ {
			hi := rand.Intn(3)
			switch hi {
			case 0:
				l.Place(l.game.NewObj(WeapSword), pt)
			case 1:
				l.Place(l.game.NewObj(ArmorLeather), pt)
			case 2:
				l.Place(l.game.NewObj(PotCure), pt)
			}
		}
	}
	l.Place(l.game.Player, math.Pt(1, 1))
	return l
}

// 80 x 80
// 15-20 rooms
// Room sizes 5-13. This includes the bounding wall, so a room with width 5
// will be 3 floor tiles across.
// Note: Thanks, fcrawl. You're a pal.
func RoomsLevel(l *Level) *Level {
	height, width, m := l.Bounds.Height(), l.Bounds.Width(), l.Map
	// When we begin, all is walls.
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			m[y][x].Feature = FeatWall
		}
	}

	nrooms := RandInt(15, 20)
	rooms := make([]math.Rectangle, 0, nrooms)
	connpts := make([]math.Point, 0, nrooms)
	log.Printf("Making %d rooms.", nrooms)

	// Find room placement.
	for i := 0; i < nrooms; i++ {
		for tries := 0; tries < (200 / nrooms); tries++ {
			rw, rh := RandInt(5, 13)|1, RandInt(5, 13)|1
			min := math.Pt(RandInt(0, width-rw), RandInt(0, height-rh))
			max := min.Add(math.Pt(rw, rh))
			newroom := math.Rect(min, max)

			log.Printf("Trying to make room %v.", newroom)
			good := true
			for _, room := range rooms {
				if newroom.Intersect(room) != math.ZeroRect {
					log.Printf("%v intersects %v -- no good.", newroom, room)
					good = false
					break
				}
			}

			if !good {
				continue
			}
			log.Printf("Room %v was good.", newroom)
			rooms = append(rooms, newroom)

			// Find a random "connection point" at some odd location in the room.
			// This is where we'll dig to from another room to make a corridor.
			connpt := math.Pt(
				RandInt(newroom.Min.X+1, newroom.Max.X-1)|1,
				RandInt(newroom.Min.Y+1, newroom.Max.Y-1)|1,
			)
			connpts = append(connpts, connpt)
		}
	}

	// Render rooms and corridors into level.
	for i, room := range rooms {
		for y := room.Min.Y + 1; y < room.Max.Y; y++ {
			for x := room.Min.X + 1; x < room.Max.X; x++ {
				m[y][x].Feature = FeatFloor
			}
		}

		if i >= 1 {
			startpt, endpt := connpts[i-1], connpts[i]

			if Coinflip() {
				for z := startpt.X; z < endpt.X; z++ {
					m[startpt.Y][z].Feature = FeatFloor
				}
				for z := startpt.Y; z < endpt.Y; z++ {
					m[z][endpt.X].Feature = FeatFloor
				}
			} else {
				for z := startpt.Y; z < endpt.Y; z++ {
					m[z][startpt.X].Feature = FeatFloor
				}
				for z := startpt.X; z < endpt.X; z++ {
					m[endpt.Y][z].Feature = FeatFloor
				}
			}
		}
	}

	startroom := rooms[RandInt(0, nrooms)]
	l.Place(l.game.Player, startroom.Center())

	return l
}

func NewDungeon(g *Game) *Level {
	return NewLevel(80, 80, g, RoomsLevel)
}
