package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
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

	if !endpos.In(obj.Map) {
		return false
	}

	return obj.Place(obj.Map, endpos)
}
