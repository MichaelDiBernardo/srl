package game

import (
	"container/heap"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type Tile struct {
	Feature *Feature
	Actor   *Obj
	Items   *Inventory
	Pos     math.Point
	Visible bool
}

type Map [][]*Tile

type Level struct {
	Map       Map
	Bounds    math.Rectangle
	game      *Game
	scheduler *Scheduler
}

// Create a level that uses the given gametory to create game objects, and which
// generated by the given generator function.
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
	default:
		panic(fmt.Sprintf("Tried to remove object of family %v", o.Spec.Family))
	}

}

// Advances the game until the player's next turn.
func (l *Level) Evolve() {
	// Advance schedule until we find the player.
	for {
		actor := l.scheduler.Next()
		if seer := actor.Seer; seer != nil {
			seer.CalcFOV()
		}
		if ai := actor.AI; ai != nil {
			ai.Act(l)
		}
		if actor.IsPlayer() {
			for _, row := range l.Map {
				for _, tile := range row {
					tile.Visible = false
				}
			}
			for _, pt := range actor.Seer.FOV() {
				l.At(pt).Visible = true
			}
			break
		}
	}
}

func (l *Level) placeActor(obj *Obj, tile *Tile) bool {
	if tile.Actor != nil {
		return false
	}

	// If this actor has been placed before, we need to clear the tile they
	// were on previously. If they haven't, we need to add them to the actor
	// list so we know who they are.
	if obj.Tile != nil {
		obj.Tile.Actor = nil
	} else {
		l.scheduler.Add(obj)
	}

	obj.Level = l
	obj.Tile = tile

	tile.Actor = obj

	return true
}

func (l *Level) placeItem(obj *Obj, tile *Tile) bool {
	return tile.Items.Add(obj)
}

func (l *Level) removeActor(obj *Obj) {
	obj.Tile.Actor = nil
	obj.Tile = nil
	obj.Level = nil
	l.scheduler.Remove(obj)
}

// Used to track which actor should be acting when.
type scheduled struct {
	actor *Obj
	delay int
}

// Shoehorning into container/heap's interface. SQ acts as the actual priority
// queue.
type SQ []*scheduled

// See https://golang.org/pkg/container/heap/#example__intHeap.
func (s SQ) Len() int {
	return len(s)
}

func (s SQ) Less(i, j int) bool {
	return s[i].delay < s[j].delay
}

func (s SQ) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s *SQ) Push(x interface{}) {
	*s = append(*s, x.(*scheduled))
}

func (s *SQ) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

type Scheduler struct {
	pq *SQ
}

func NewScheduler() *Scheduler {
	pq := make(SQ, 0)
	return &Scheduler{pq: &pq}
}

func (s *Scheduler) Len() int {
	return s.pq.Len()
}

// Add an actor to the schedule.
func (s *Scheduler) Add(actor *Obj) {
	delay := getdelay(actor.Sheet.Speed())

	// Player should always get the first turn.
	if actor.IsPlayer() {
		delay = 0
	}

	entry := &scheduled{actor: actor, delay: delay}
	heap.Push(s.pq, entry)
}

// Picks the next actor to act, and moves time forward for everyone else.
func (s *Scheduler) Next() *Obj {
	entry := heap.Pop(s.pq).(*scheduled)
	for _, e := range *(s.pq) {
		e.delay -= entry.delay
	}
	heap.Init(s.pq)

	actor := entry.actor
	entry.delay = getdelay(actor.Sheet.Speed())
	heap.Push(s.pq, entry)

	return actor
}

// Removes an actor from the scheduler.
func (s *Scheduler) Remove(actor *Obj) {
	index := -1
	for i, e := range *(s.pq) {
		if e.actor == actor {
			index = i
			break
		}
	}

	if index == -1 {
		panic("Tried to remove actor but wasn't in list.")
	}

	*(s.pq) = append((*(s.pq))[:index], (*(s.pq))[index+1:]...)
	heap.Init(s.pq)
}

// Given the speed of an actor, this will tell you how much delay to add after
// each of its turns.
func getdelay(spd int) int {
	switch spd {
	case 1:
		return 150
	case 2:
		return 100
	case 3:
		return 75
	case 4:
		return 50
	default:
		panic(fmt.Sprintf("Spd %d does not have a delay", spd))
	}
}
