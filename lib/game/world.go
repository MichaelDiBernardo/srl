package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type World struct {
	Player *Player
}

func NewWorld() *World {
	return &World{Player: NewPlayer()}
}

func (w *World) Handle(e Command) {
	player := w.Player
	switch e {
	case CommandMoveN:
		player.Pos = player.Pos.Add(math.Point{0, -1})
	case CommandMoveS:
		player.Pos = player.Pos.Add(math.Point{0, 1})
	case CommandMoveW:
		player.Pos = player.Pos.Add(math.Point{-1, 0})
	case CommandMoveE:
		player.Pos = player.Pos.Add(math.Point{1, 0})
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
