package game

import (
	"errors"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type Mover interface {
	Objgetter
	Move(dir math.Point) error
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

var (
	ErrMove0Dir        = errors.New("Move0Dir")
	ErrMoveTooFar      = errors.New("MoveTooFar")
	ErrMoveHit         = errors.New("MoveHit")
	ErrMoveBlocked     = errors.New("MoveBlocked")
	ErrMoveOutOfBounds = errors.New("MoveOutOfBounds")
	ErrMoveSwapFailed  = errors.New("MoveSwapFailed")
	ErrMoveOpenedDoor  = errors.New("MoveOpenedDoor")
)

// Try to move the actor. Return err describing what happened if the move
// fails.
func (p *ActorMover) Move(dir math.Point) error {
	if dist := math.ChebyDist(math.Origin, dir); dist > 1 {
		return ErrMoveTooFar
	} else if dist == 0 {
		return ErrMove0Dir
	}

	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return ErrMoveOutOfBounds
	}

	endtile := obj.Level.At(endpos)
	if other := endtile.Actor; other != nil {
		if opposing := obj.IsPlayer() != other.IsPlayer(); opposing {
			p.obj.Fighter.Hit(other.Fighter)
			return ErrMoveHit
		} else {
			// Traveling monsters should swap with one another, but it's kind
			// of a pain.
			if OneIn(2) {
				obj.Level.SwapActors(obj, other)
				return nil
			}
			return ErrMoveSwapFailed
		}
	}
	if endtile.Feature == FeatClosedDoor {
		endtile.Feature = FeatOpenDoor
		return ErrMoveOpenedDoor
	}

	moved := obj.Level.Place(obj, endpos)
	if moved {
		if items := endtile.Items; !items.Empty() && obj.IsPlayer() && !obj.Sheet.Blind() {
			var msg string
			topname, n := items.Top().Spec.Name, items.Len()
			if n == 1 {
				msg = fmt.Sprintf("%v sees %v here.", obj.Spec.Name, topname)
			} else {
				msg = fmt.Sprintf("%v sees %v and %d other items here.", obj.Spec.Name, topname, n-1)
			}
			obj.Game.Events.Message(msg)
		}
		return nil
	}
	return ErrMoveBlocked
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
