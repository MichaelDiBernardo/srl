package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type World struct {
	Player *Obj
	Map    Map
}

func NewWorld() *World {
	wmap := NewMap(80, 24, SquareMap)
	player := NewObj(Traits{Mover: NewActorMover})
	player.Place(wmap, math.Pt(1, 1))
	return &World{Player: player, Map: wmap}
}

func (w *World) Handle(e Command) {
	switch e {
	case CommandMoveN:
		w.Player.Mover.Move(math.Pt(0, -1))
	case CommandMoveS:
		w.Player.Mover.Move(math.Pt(0, 1))
	case CommandMoveW:
		w.Player.Mover.Move(math.Pt(-1, 0))
	case CommandMoveE:
		w.Player.Mover.Move(math.Pt(1, 0))
	}
}

type Command int

const (
	_ Command = iota
	CommandQuit
	CommandMoveN
	CommandMoveNE
	CommandMoveE
	CommandMoveSE
	CommandMoveS
	CommandMoveSW
	CommandMoveW
	CommandMoveNW
)
