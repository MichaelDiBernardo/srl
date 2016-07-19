package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

type Player struct {
	Tile *Tile
	Map  Map
}

func (p *Player) GetPos() math.Point {
	return p.Tile.Pos
}

// Try to move the player. Return false if the player couldn't move.
func (p *Player) Move(dir math.Point) bool {
	beginpos := p.GetPos()
	endpos := p.GetPos().Add(dir)

	// TODO: Bounds check
	// TODO: Map needs a type with methods for these accesses.
	begintile := p.Map[beginpos.X][beginpos.Y]
	endtile := p.Map[endpos.X][endpos.Y]

	if endtile.Feature.Solid {
		return false
	}

	begintile.Actor = nil
	endtile.Actor = p
	p.Tile = endtile
	return true
}
