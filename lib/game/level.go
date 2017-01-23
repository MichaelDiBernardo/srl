package game

import (
	"container/heap"
	"fmt"
	num "math"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

type FeatureType string

type Feature struct {
	Type   FeatureType
	Solid  bool
	Opaque bool
}

func (f *Feature) String() string {
	return string(f.Type)
}

type Tile struct {
	Feature *Feature
	Actor   *Obj
	Items   *Inventory
	Pos     math.Point
	Visible bool
	Seen    bool
	Scent   int
}

func (t *Tile) String() string {
	return fmt.Sprintf("<%s:%v>", t.Feature, t.Pos)
}

type Map [][]*Tile

type Level struct {
	Map       Map
	Bounds    math.Rectangle
	game      *Game
	scheduler *Scheduler
}

// Create a level that uses the given game to create objects, generated by the
// given generator function.
func NewLevel(width, height int, game *Game, gen func(*Level) *Level) *Level {
	newmap := Map{}

	// Init all tiles.
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		for x := 0; x < width; x++ {
			row[x] = &Tile{
				Pos:     math.Pt(x, y),
				Feature: FeatFloor,
				Items:   NewInventory(),
			}
		}
		newmap = append(newmap, row)
	}
	level := &Level{
		Map:       newmap,
		Bounds:    math.Rect(math.Origin, math.Pt(width, height)),
		game:      game,
		scheduler: NewScheduler(),
	}
	level = gen(level)

	// Init all the actors brains, now that they have a place on the map.
	level.scheduler.EachActor(func(o *Obj) {
		if o.AI != nil {
			o.AI.Init()
		}
	})
	return level
}

func (l *Level) At(p math.Point) *Tile {
	return l.Map[p.Y][p.X]
}

// Returns a list of the tiles "around" p (but not including p.)
func (l *Level) Around(loc math.Point) []*Tile {
	adj := l.Bounds.Clip(math.Adj(loc))

	neighbors := make([]*Tile, 0, len(adj))

	for _, p := range adj {
		neighbors = append(neighbors, l.At(p))
	}
	return neighbors
}

func (l *Level) HasPoint(p math.Point) bool {
	return l.Bounds.HasPoint(p)
}

// Place `o` on the tile at `p`. Returns false if this is impossible (e.g.
// trying to put something on a solid square.)
// This will remove `o` from any tile on any map it was previously on.
func (l *Level) Place(o *Obj, p math.Point) bool {
	t := l.At(p)
	if t.Feature.Solid {
		return false
	}

	switch o.Spec.Family {
	case FamActor:
		return l.placeActor(o, t)
	case FamItem:
		return l.placeItem(o, t)
	default:
		panic(fmt.Sprintf("Tried to place object of family %v", o.Spec.Family))
	}
}

// Removes 'o' from the level.
func (l *Level) Remove(o *Obj) {
	switch o.Spec.Family {
	case FamActor:
		l.removeActor(o)
		l.scheduler.Remove(o)
	default:
		panic(fmt.Sprintf("Tried to remove object of family %v", o.Spec.Family))
	}

}

// Advances the game until the player's next turn.
func (l *Level) Evolve() {
	// Advance schedule until we find the player.
	for {
		actor := l.scheduler.Next()
		// TODO: This may need to be done on _every_ actor for each turn taken.
		actor.Ticker.Tick(l.scheduler.delay)

		if ai := actor.AI; ai != nil {
			ai.Act()
		}
		if actor.IsPlayer() {
			break
		}
	}
}

func (l *Level) RandomClearTile() *Tile {
	r := l.Bounds
	for tries := 0; tries < 100; tries++ {
		loc := math.Pt(
			RandInt(r.Min.X, r.Max.X),
			RandInt(r.Min.Y, r.Max.Y),
		)
		tile := l.At(loc)
		if !tile.Feature.Solid {
			return tile
		}
	}
	return nil
}

func (l *Level) SwapActors(x *Obj, y *Obj) {
	tx, ty := x.Tile, y.Tile
	l.removeActor(y)
	l.placeActor(x, ty)
	l.placeActor(y, tx)
}

// Update what the player can see.
func (l *Level) UpdateVis(fov Field) {
	for _, row := range l.Map {
		for _, tile := range row {
			tile.Visible = false
		}
	}

	for _, pt := range fov {
		tile := l.At(pt)
		tile.Visible = true
		tile.Seen = true
	}
}

// Update the player's smell on the map.
func (l *Level) UpdateScent(scent Field) {
	pos, turns := l.game.Player.Pos(), l.game.Turns

	for _, pt := range scent {
		tile := l.At(pt)
		tile.Scent = turns*ScentRadius - math.ChebyDist(pos, pt)
	}
}

// Update the given field on the map.

// Return type of FindPath. This is what you follow to travel the found route.
type Path []math.Point

// The distance map in Dijkstra's algorithm.
type dists map[math.Point]int

// Return basically INFINITY if we haven't seen p before.
func (d dists) get(p math.Point) int {
	dist, ok := d[p]
	if ok {
		return dist
	}
	return num.MaxInt32
}

// An unvisited position on the map. We hold a reference to the distances dict
// for distance lookup instead of trying to maintain it all in one place.
type unvisited struct {
	p    math.Point
	dist dists
}

// Look up the currently calculated distance for this point.
func (u *unvisited) d() int {
	return u.dist.get(u.p)
}

// The PQ for dijkstra pathfinding.
type dijkstraQ []*unvisited

// Implementing heap.Interface
func (d dijkstraQ) Len() int {
	return len(d)
}

func (d dijkstraQ) Less(i, j int) bool {
	return d[i].d() < d[j].d()
}

func (d dijkstraQ) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d *dijkstraQ) Push(x interface{}) {
	*d = append(*d, x.(*unvisited))
}

func (d *dijkstraQ) Pop() interface{} {
	old := *d
	n := len(old)
	x := old[n-1]
	*d = old[0 : n-1]
	return x
}

// Finds a reasonably direct path between start and dest in this level. If no
// path could be found, 'path' will be zero and 'ok' will be false. This uses
// Dijkstra's algorithm.
func (l *Level) FindPath(start, end math.Point, cost func(*Level, math.Point) int) (path Path, ok bool) {
	path = make(Path, 0, 100)
	if !start.In(l) || !end.In(l) || l.At(start).Feature.Solid || l.At(end).Feature.Solid {
		return path, false
	}

	if start == end {
		return path, true
	}

	// PQ
	dist := dists{start: 0}
	prev := map[math.Point]math.Point{}
	todo := &dijkstraQ{&unvisited{p: start, dist: dist}}
	heap.Init(todo)

	for todo.Len() != 0 {
		cur := heap.Pop(todo).(*unvisited).p
		curdist := dist[cur]
		if cur == end {
			break
		}

		// Get neighbouring points.
		neighbors := l.Around(cur)

		// Figure out which neighbours we should even bother checking (e.g.
		// walls are a bad idea at the moment. We don't have Kemenrauko.
		for _, n := range neighbors {
			if !patheligible(n) {
				continue
			}
			npos := n.Pos
			d := dist.get(npos)

			altdist := curdist + cost(l, npos)
			if altdist < d {
				dist[npos] = altdist
				prev[npos] = cur
				heap.Push(todo, &unvisited{p: npos, dist: dist})
			}
		}
	}

	// Trace path, which will give us the route in reverse.
	pathcur := end
	for {
		next, ok := prev[pathcur]
		if !ok {
			break
		}
		path = append(path, pathcur)
		pathcur = next
	}

	// Put the path in order.
	for i, l := 0, len(path); i < l/2; i++ {
		path[i], path[l-i-1] = path[l-i-1], path[i]
	}

	// We didn't find any path. Let the caller know via 'ok', since this
	// distinguishes the situation where you tried to pathfind to yourself.
	if len(path) == 0 {
		return path, false
	}
	return path, true
}

// Returns the "cost" of moving onto 'loc' in level l.
func PathCost(l *Level, loc math.Point) int {
	switch l.At(loc).Feature {
	case FeatClosedDoor:
		return 2
	default:
		return 1
	}
}

func patheligible(t *Tile) bool {
	return t.Feature != FeatWall
}

func (l *Level) placeActor(obj *Obj, tile *Tile) bool {
	if tile.Actor != nil {
		return false
	}

	// If this actor has been placed before, we need to clear the tile they
	// were on previously. If they haven't, we need to add them to the actor
	// list so we know who they are.
	if obj.Tile != nil && obj.Level == l {
		obj.Tile.Actor = nil
	} else {
		l.scheduler.Add(obj)
	}

	obj.Level = l
	obj.Tile = tile

	tile.Actor = obj

	// Refresh any fields for the actor here.
	if ticker := obj.Ticker; ticker != nil {
		ticker.Tick(0)
	}

	return true
}

func (l *Level) placeItem(obj *Obj, tile *Tile) bool {
	return tile.Items.Add(obj)
}

// Removes actor from the board. Unlike placeActor, this does NOT change the
// actor's place in the scheduler; this is so other Level methods (like
// SwapActors) can lift and replace some dudes without rescheduling them.
func (l *Level) removeActor(obj *Obj) {
	obj.Tile.Actor = nil
	obj.Tile = nil
	obj.Level = nil
}
