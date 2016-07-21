package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type Mover interface {
	Move(dir math.Point) bool
}

type PlayerMover struct {
	Obj *Obj
}

func NewPlayerMover(obj *Obj) Mover {
	return &PlayerMover{Obj: obj}
}

// Try to move the player. Return false if the player couldn't move.
func (p *PlayerMover) Move(dir math.Point) bool {
	obj := p.Obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Map) {
		return false
	}

	return obj.Place(obj.Map, endpos)
}
