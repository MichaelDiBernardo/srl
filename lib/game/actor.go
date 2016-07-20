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

	if !endpos.In(p.Map) {
		return false
	}

	begintile := p.Map.At(beginpos)
	endtile := p.Map.At(endpos)

	if endtile.Feature.Solid {
		return false
	}

	begintile.Actor = nil
	endtile.Actor = p
	p.Tile = endtile
	return true
}
