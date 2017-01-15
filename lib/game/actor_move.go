package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type Mover interface {
	Objgetter
	Move(dir math.Point) bool
	Ascend() bool
	Descend() bool
}

// A universally-applicable mover for actors.
type ActorMover struct {
	Trait
}

// Constructor for actor movers.
func NewActorMover(obj *Obj) Mover {
	return &ActorMover{Trait: Trait{obj: obj}}
}

// Try to move the actor. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	if math.ChebyDist(math.Origin, dir) > 1 {
		panic(fmt.Sprintf("Tried to Move(%s)", dir))
	}
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	endtile := obj.Level.At(endpos)
	if other := endtile.Actor; other != nil {
		if opposing := obj.IsPlayer() != other.IsPlayer(); opposing {
			p.obj.Fighter.Hit(other.Fighter)
			return false
		} else {
			// Traveling monsters should swap with one another.
			obj.Level.SwapActors(obj, other)
			return true
		}
	}
	if endtile.Feature == FeatClosedDoor {
		endtile.Feature = FeatOpenDoor
		return false
	}

	moved := obj.Level.Place(obj, endpos)
	if moved {
		if items := endtile.Items; !items.Empty() && obj.IsPlayer() {
			var msg string
			topname, n := items.Top().Spec.Name, items.Len()
			if n == 1 {
				msg = fmt.Sprintf("%v sees %v here.", obj.Spec.Name, topname)
			} else {
				msg = fmt.Sprintf("%v sees %v and %d other items here.", obj.Spec.Name, topname, n-1)
			}
			obj.Game.Events.Message(msg)
		}
	}
	return moved
}

// Try to go up stairs. If the current tile is not an upstair, return false.
func (p *ActorMover) Ascend() bool {
	if tile := p.obj.Tile; tile.Feature != FeatStairsUp {
		return false
	}
	p.obj.Game.ChangeFloor(1)
	return true
}

// Try to go down stairs. If the current tile is not an downstair, return false.
func (p *ActorMover) Descend() bool {
	if tile := p.obj.Tile; tile.Feature != FeatStairsDown {
		return false
	}
	p.obj.Game.ChangeFloor(-1)
	return true
}
