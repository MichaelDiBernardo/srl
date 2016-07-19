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

func (w *World) Handle(e Event) {
	player := w.Player
	switch e {
	case EMoveN:
		player.Pos = player.Pos.Add(math.Point{0, -1})
	case EMoveS:
		player.Pos = player.Pos.Add(math.Point{0, 1})
	case EMoveW:
		player.Pos = player.Pos.Add(math.Point{-1, 0})
	case EMoveE:
		player.Pos = player.Pos.Add(math.Point{1, 0})
	}
}
