package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Backend for a single game.
type Game struct {
	Player *Obj
	Level  *Level
}

func NewGame() *Game {
	level := NewLevel(80, 24, TestLevel)
	player := NewPlayer()
	level.Place(player, math.Pt(1, 1))
	return &Game{Player: player, Level: level}
}

func (w *Game) Handle(e Command) {
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
	w.Level.Evolve()
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
