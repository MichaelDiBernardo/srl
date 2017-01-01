package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"math/rand"
)

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Act(l *Level) bool
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	Trait
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{Trait: Trait{obj: obj}}
}

// Move in any of the 8 directions with uniform chance. Does not take walls
// etc. in account so this will happily try to bump into things.
func (ai *RandomAI) Act(l *Level) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	dir := math.Pt(x, y)
	if dir == math.Origin {
		return ai.Act(l)
	}
	return ai.obj.Mover.Move(dir)
}
