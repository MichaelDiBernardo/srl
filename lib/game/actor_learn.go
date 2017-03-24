package game

import (
	"fmt"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Tracks xp gained by monsters and items seen and killed.
type Learner interface {
	// Call this when a specific instance of an item or monster has been seen.
	LearnSight(obj *Obj)
	// Call this when a specific instance of a monster has been seen.
	LearnKill(mon *Obj)
	// How much XP does this actor have?
	XP() int
}

// Used by the player to track what they have seen, and how much XP they have.
type ActorLearner struct {
	seen   map[Species]int
	killed map[Species]int
	xp     int
}

func NewActorLearner(obj *Obj) Learner {
	return &ActorLearner{
		seen:   map[Species]int{},
		killed: map[Species]int{},
	}
}

func (l *ActorLearner) XP() int {
	return l.xp
}

func (l *ActorLearner) LearnKill(mon *Obj) {
	if genus := mon.Spec.Genus; genus != GenMonster {
		panic(fmt.Sprintf("Obj *v with genus %v is not monster.", mon, genus))
	}

	s := mon.Spec.Species
	n := l.killed[s]

	xpgain := monxp(mon, n)
	l.xp += xpgain
	l.killed[s]++
}

func (l *ActorLearner) LearnSight(obj *Obj) {
	if obj.Seen {
		return
	}

	s, xp := obj.Spec.Species, 0
	n := l.seen[s]

	switch genus := obj.Spec.Genus; genus {
	case GenMonster:
		xp = monxp(obj, n)
	case GenEquipment, GenConsumable:
		xp = itemxp(obj, n)
	default:
		panic(fmt.Sprintf("Obj *v with genus %v is not xpable on sight.", obj, genus))
	}

	l.xp += xp
	l.seen[s]++
}

func monxp(mon *Obj, n int) int {
	return xpdecay(xpfor(mon), n)
}

func itemxp(item *Obj, n int) int {
	if n > 0 {
		return 0
	}
	return xpfor(item)
}

func xpfor(obj *Obj) int {
	depth := obj.Spec.Gen.First()
	return depth * 10
}

func xpdecay(xp, n int) int {
	return math.Max(1, xp/(n+1))
}
