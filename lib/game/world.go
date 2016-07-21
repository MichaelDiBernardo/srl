package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type Tile struct {
	Feature *Feature
	Actor   *Obj
	Pos     math.Point
}

type Map [][]*Tile

func NewMap(width, height int) Map {
	newmap := Map{}
	for y := 0; y < height; y++ {
		row := make([]*Tile, width, width)
		newmap = append(newmap, row)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			feature := FeatFloor
			if x == 0 || y == 0 || y == height-1 || x == width-1 {
				feature = FeatWall
			}
			newmap[y][x] = &Tile{Pos: math.Pt(x, y), Feature: feature}
		}
	}
	return newmap
}

func (m Map) Width() int {
	return len(m[0])
}

func (m Map) Height() int {
	return len(m)
}

func (m Map) At(p math.Point) *Tile {
	return m[p.Y][p.X]
}

func (m Map) HasPoint(p math.Point) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < m.Width() && p.Y < m.Height()
}

type World struct {
	Player *Obj
	Map    Map
}

func NewWorld() *World {
	wmap := NewMap(80, 24)
	player := NewObj(Traits{Mover: NewPlayerMover})
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
