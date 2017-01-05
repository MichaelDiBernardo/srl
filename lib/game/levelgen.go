package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
)

// Generates a simple square test level suitable for unit tests.
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
	l.Place(l.game.Player, math.Pt(2, 2))
	return l
}

// This implements lynn's algorithm for laying out angband-ish rooms
// without a lot of fuss.
//
// The idea is to randomly place rooms, and then draw L-shaped corridors to
// connect them. The connection points are randomly-selected "joints" placed
// one-per-room.
//
// We do this by:
// a) Creating and randomly placing an odd-sized room where it fits. We don't
//    allow new rooms to overlap existing rooms or corridors.
// b) Selecting a random, odd-aligned joint in this new room.
// c) If there is at least one predecessor room, dig an L-shaped path between
//    the joints. At this stage, corridors are allowed to blast through
//    intervening rooms.
//
// See http://i.imgur.com/WhmnByV.png for a terse graph-ical explanation.
func LynnRoomsLevel(l *Level) *Level {
	// When we begin, all is walls.
	m := l.Map
	fillmap(m, l.Bounds, FeatWall)

	// We'll attempt to create this many rooms, but may fall short if we run
	// into intractable placement problems.
	maxrooms := RandInt(10, 20)

	// The rooms we've placed so far, represented as Rects.
	rooms := make([]math.Rectangle, 0, maxrooms)
	// The ith joint in 'joints' is the joint in the ith room of 'rooms'.
	joints := make([]math.Point, 0, maxrooms)
	// All of the points that currently belong to a corridor.
	paths := make([]math.Point, 0)

	log.Printf("Making %d rooms.", maxrooms)

	// The actual room index and # of rooms.
	ri, nrooms := 0, 0
	for r := 0; r < maxrooms; r++ {
		placed := placeroom(l, rooms, paths)
		if placed == math.ZeroRect {
			// We tried hard but couldn't do it. Try again from beginning.
			continue
		}

		// Record and draw the room on the map. We want it drawn before we try
		// to dig out of it to make door placement easier.
		rooms = append(rooms, placed)
		fillmap(m, placed, FeatFloor)

		// Find a joint for this room, and if we're far along enough, try to
		// join it to the previous.
		joints = append(joints, makejoint(placed))
		if ri > 0 {
			path := dig(joints[ri], joints[ri-1])
			drawpath(l, path, rooms)
			for _, pt := range path {
				paths = append(paths, pt)
			}
		}
		ri++
		nrooms++
	}
	path := dig(joints[0], joints[nrooms-1])
	drawpath(l, path, rooms)
	// Don't need to add to paths anymore since we're done placing rooms.

	startroom := rooms[RandInt(0, nrooms)]
	l.Place(l.game.Player, startroom.Center())

	placemonsters(l, startroom, rooms)
	placeitems(l, rooms)

	log.Printf("Made %d/%d rooms.", nrooms, maxrooms)
	return l
}

// Fills the tiles in the given area with the given feature.
func fillmap(m Map, area math.Rectangle, f *Feature) {
	min, max := area.Min, area.Max
	for y := min.Y; y < max.Y; y++ {
		for x := min.X; x < max.X; x++ {
			m[y][x].Feature = f
		}
	}
}

// Makes a random room within the confines of the given level.
func randroom(l *Level) math.Rectangle {
	width, height := l.Bounds.Width(), l.Bounds.Height()

	// Clamp room to odd widths and heights -- this is an aesthetic preference.
	// The left side of the range is an even number to make all of the possible
	// odd values equally probably ([4..5]=5, [6..7]=7 etc.)
	rw, rh := RandInt(4, 13)|1, RandInt(4, 13)|1
	min := math.Pt(RandInt(1, width-rw-3), RandInt(1, height-rh-3))
	max := min.Add(math.Pt(rw, rh))
	return math.Rect(min, max)
}

// Tries to generate a room that won't intersect with any other rooms or an
// existing corridor.
func placeroom(l *Level, rooms []math.Rectangle, paths []math.Point) math.Rectangle {
	log.Printf("Placeroom:")
	for tries := 0; tries < 15; tries++ {
		newroom := randroom(l)
		log.Printf("\tTrying to place room candidate: %v.", newroom)
		if fits(newroom, rooms, paths) {
			return newroom
		}
	}
	log.Print("\tCouldn't place any candidates; giving up.")
	return math.ZeroRect
}

// Can 'room' be placed in a level with 'rooms' and 'paths' without
// intersecting any of them?
func fits(newroom math.Rectangle, rooms []math.Rectangle, paths []math.Point) bool {
	nrbounds := math.Rect(
		newroom.Min.Add(math.Pt(-1, -1)),
		newroom.Max.Add(math.Pt(1, 1)),
	)

	for _, room := range rooms {
		// When checking for intersection, we'll use a room boundary 1
		// larger than the floorspace of the actual room.
		// This is so we don't get rooms that are directly adjacent to one
		// another (no "frankenrooms".)
		if nrbounds.Intersect(room) != math.ZeroRect {
			log.Printf("%v intersects %v -- no good.", newroom, room)
			return false
		}
	}

	for _, pt := range paths {
		if nrbounds.HasPoint(pt) {
			return false
		}
	}

	return true
}

// Given two joints, this will return a path that joins them.
func dig(startpt, endpt math.Point) []math.Point {
	log.Printf("Connecting joint %v to %v", startpt, endpt)

	var start, end, incr int
	path := make([]math.Point, 0)

	if Coinflip() {
		start, end, incr = drange(startpt.X, endpt.X, true)
		for z := start; z != end; z += incr {
			pt := math.Pt(z, startpt.Y)
			path = append(path, pt)
			log.Printf("\t%v", pt)
		}
		start, end, incr = drange(startpt.Y, endpt.Y, true)
		for z := start; z != end; z += incr {
			pt := math.Pt(endpt.X, z)
			path = append(path, pt)
			log.Printf("\t%v", pt)
		}
	} else {
		start, end, incr = drange(startpt.Y, endpt.Y, true)
		for z := start; z != end; z += incr {
			pt := math.Pt(startpt.X, z)
			path = append(path, pt)
			log.Printf("\t%v", pt)
		}
		start, end, incr = drange(startpt.X, endpt.X, true)
		for z := start; z != end; z += incr {
			pt := math.Pt(z, endpt.Y)
			path = append(path, pt)
			log.Printf("\t%v", pt)
		}
	}

	return path
}

// Finds an odd-aligned location to serve as a joint in an l-shaped path
// connecting this to another room.
func makejoint(room math.Rectangle) math.Point {
	return math.Pt(
		RandInt(room.Min.X, room.Max.X)|1,
		RandInt(room.Min.Y, room.Max.Y)|1,
	)
}

// Draws the path from startroom to endroom. Also places closed doors at egress
// of each room that is intersected along the way.
func drawpath(l *Level, path []math.Point, rooms []math.Rectangle) {
	// Predicate that tells us if we should place a door.
	placedoor := func(i int, pt math.Point) bool {
		if i == 0 || i >= len(path)-1 {
			return false
		}
		prevpt, nextpt := path[i-1], path[i+1]
		for _, room := range rooms {
			egress := !room.HasPoint(pt) && room.HasPoint(prevpt)
			ingress := !room.HasPoint(pt) && room.HasPoint(nextpt)
			if egress || ingress {
				return true
			}
		}
		return false
	}
	for i, pt := range path {
		if tile := l.At(pt); placedoor(i, pt) {
			log.Printf("Placed door at %v", pt)
			tile.Feature = FeatClosedDoor
		} else {
			tile.Feature = FeatFloor
		}
	}
}

// Generates and places monsters in any room except the starting room.
func placemonsters(l *Level, startroom math.Rectangle, rooms []math.Rectangle) {
	g := l.game

	mongroups := GenMonsters(20, g.Depth, 2, g)

	for _, group := range mongroups {
		for tries := 0; tries < 50; tries++ {
			room := rooms[RandInt(0, len(rooms))]
			if room == startroom {
				continue
			}

			// TODO: Don't use another loop, instead try to place as a group.
			for _, mon := range group {
				for moretries := 0; moretries < 50; moretries++ {
					loc := math.Pt(
						RandInt(room.Min.X, room.Max.X),
						RandInt(room.Min.Y, room.Max.Y),
					)

					if l.Place(mon, loc) {
						break
					}
				}
			}
		}
	}
}

// Generates and places a bunch of items in any room.
func placeitems(l *Level, rooms []math.Rectangle) {
	nitems := 40
	g := l.game

	for i := 0; i < nitems; i++ {
		var spec *Spec
		switch RandInt(0, 3) {
		case 0:
			spec = WeapSword
		case 1:
			spec = ArmorLeather
		case 2:
			spec = PotCure
		}

		item := g.NewObj(spec)
		for tries := 0; tries < 10; tries++ {
			room := rooms[RandInt(0, len(rooms))]
			loc := math.Pt(
				RandInt(room.Min.X, room.Max.X),
				RandInt(room.Min.Y, room.Max.Y),
			)
			if l.Place(item, loc) {
				break
			}
		}
	}
}

// Given x and y, this will return 'start', 'end', and 'iter' that you can use
// in a for loop to iterate in an ordered sequence from x to y. Set 'incl' to
// true if you want both x and y to be included in the iteration; otherwise it
// will end just before y.
//
// Use it like:
// start, end, incr, i := drange(10, 4, true)
// for i := start; i != end; i += incr {
//    dostuff()
// }
func drange(x, y int, incl bool) (start, end, incr int) {
	if x < y {
		start, end, incr = x, y, 1
	} else {
		start, end, incr = x, y, -1
	}

	if incl {
		end += incr
	}
	return start, end, incr
}

// Generates a new dungeon level.
func NewDungeon(g *Game) *Level {
	return NewLevel(80, 80, g, LynnRoomsLevel)
}
