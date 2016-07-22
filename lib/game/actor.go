package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"math/rand"
)

// A thing that can move given a specific direction.
type Mover interface {
	Move(dir math.Point) bool
}

type ActorMover struct {
	Obj *Obj
}

func NewActorMover(obj *Obj) Mover {
	return &ActorMover{Obj: obj}
}

// Try to move the player. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.Obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	return obj.Place(obj.Level, endpos)
}

// A thing that can move given a specific direction.
type AI interface {
	Act(w World) bool
}

type RandomAI struct {
	Obj *Obj
}

func NewRandomAI(obj *Obj) AI {
	return &RandomAI{Obj: obj}
}

func (ai *RandomAI) Act(w World) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1

	return ai.Obj.Mover.Move(math.Pt(x, y))
}
