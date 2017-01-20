package game

import (
	"container/heap"
	"fmt"
)

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

// Schedules which actors should act when.
type Scheduler struct {
	// The priority queue that orders the actors.
	pq *SQ
	// A tiny variance in initial delay to enforce ordering (see later
	// comments.)
	bump int
	// How much total delay has been processed on this scheduler.
	delay int
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
	// Attempting to re-add anyone to the scheduler won't refresh their current
	// delay. We just ignore the request.
	for _, e := range *(s.pq) {
		if e.actor == actor {
			return
		}
	}

	// In Next(), we adjust everyone's delay and then call Init() to re-heapify
	// pq. This isn't a stable "sort", so actors of equal speed end up randomly
	// taking their turns in no particular order if they all start the game at
	// the same delay. So, we add a tiny 'bump' to their delay which will
	// synchronize their ordering in subsequent turns, since they only ever
	// lose delay in much larger increments.
	delay := GetDelay(actor.Sheet.Speed()) + s.bump
	s.bump++

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
	s.delay += entry.delay
	heap.Init(s.pq)

	actor := entry.actor
	entry.delay = GetDelay(actor.Sheet.Speed())
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

	heap.Remove(s.pq, index)
}

// Given the speed of an actor, this will tell you how much delay to add after
// each of its turns.
func GetDelay(spd int) int {
	switch spd {
	case 1:
		return 1500
	case 2:
		return 1000
	case 3:
		return 750
	case 4:
		return 500
	default:
		panic(fmt.Sprintf("Spd %d does not have a delay", spd))
	}
}
