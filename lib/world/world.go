package world

import (
	"github.com/MichaelDiBernardo/srl/lib/actor"
	"github.com/MichaelDiBernardo/srl/lib/event"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type World struct {
	Player *actor.Player
}

func New() *World {
	return &World{Player: actor.NewPlayer()}
}

func (w *World) Handle(e event.Event) {
	player := w.Player
	switch e {
	case event.MoveN:
		player.Pos = player.Pos.Add(math.Point{0, -1})
	case event.MoveS:
		player.Pos = player.Pos.Add(math.Point{0, 1})
	case event.MoveW:
		player.Pos = player.Pos.Add(math.Point{-1, 0})
	case event.MoveE:
		player.Pos = player.Pos.Add(math.Point{1, 0})
	}
}
