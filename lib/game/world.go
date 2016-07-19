package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type Tile struct {
	Feature *Feature
	Actor   *Player
	Pos     math.Point
}

type Map [80][24]*Tile

type World struct {
	Player *Player
	Map    Map
}

func NewWorld() *World {
	wmap := Map{}
	for i := 0; i < 80; i++ {
		for j := 0; j < 24; j++ {
			feature := FeatFloor
			if i == 0 || j == 0 || i == 79 || j == 23 {
				feature = FeatWall
			}
			wmap[i][j] = &Tile{Pos: math.Pt(i, j), Feature: feature}
		}
	}
	player := &Player{
		Map:  wmap,
		Tile: wmap[1][1],
	}
	wmap[1][1].Actor = player
	return &World{Player: player, Map: wmap}
}

func (w *World) Handle(e Command) {
	switch e {
	case CommandMoveN:
		w.Player.Move(math.Pt(0, -1))
	case CommandMoveS:
		w.Player.Move(math.Pt(0, 1))
	case CommandMoveW:
		w.Player.Move(math.Pt(-1, 0))
	case CommandMoveE:
		w.Player.Move(math.Pt(1, 0))
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
