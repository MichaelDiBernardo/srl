package game

import (
	"errors"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type Mover interface {
	Objgetter
	Move(dir math.Point) (bool, error)
	Rest()
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
	ErrTooScaredToHit  = errors.New("TooScaredToHit")
	ErrMoveBlocked     = errors.New("MoveBlocked")
	ErrMoveOutOfBounds = errors.New("MoveOutOfBounds")
	ErrMoveSwapFailed  = errors.New("MoveSwapFailed")
	ErrMoveOpenedDoor  = errors.New("MoveOpenedDoor")
)

// Try to move the actor. Return err describing what happened if the actor
// could not physically move on the map; there may have been another event that
// had to happen instead that still requires a turn to pass (e.g. a door was
// opened, an attack was initiated.) If a turn should pass, the boolean return
// value from this function will be true. This idea of a 'turn passing' is in
// the context of the player's turn -- calling code can decide whether or not a
// monster's turn should be over, regardless of this value.
func (p *ActorMover) Move(dir math.Point) (bool, error) {
	obj := p.obj
	conf := obj.Sheet.Confused()

	// Validate weirdness.
	// Rest() may tell Move() to move 0,0 if the actor is confused, so that's
	// why we check 'conf' in the second clause.
	if math.ChebyDist(math.Origin, dir) > 1 {
		return false, ErrMoveTooFar
	} else if dir == math.Origin && !conf {
		return false, ErrMove0Dir
	}

	if conf {
		if OneIn(2) {
			dir = confusedir(dir)
			obj.Game.Events.Message(fmt.Sprintf("%v moves the wrong way.", obj.Spec.Name))
		} else if dir == math.Origin {
			// We're confused and we tried to pass a turn, but we didn't have a
			// direction chosen for us. This consumes a turn, but we're still
			// not going to move anywhere.
			return true, ErrMove0Dir
		}
	}

	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return conf || false, ErrMoveOutOfBounds
	}

	endtile := obj.Level.At(endpos)
	if other := endtile.Actor; other != nil {
		if opposing := obj.IsPlayer() != other.IsPlayer(); opposing {
			if obj.Sheet.Afraid() {
				msg := fmt.Sprintf("%s is too afraid to attack %s!", obj.Spec.Name, other.Spec.Name)
				obj.Game.Events.Message(msg)
				return false, ErrTooScaredToHit
			}
			p.obj.Fighter.Hit(other.Fighter)
			return true, ErrMoveHit
		} else {
			// Traveling monsters should swap with one another, but it's kind
			// of a pain.
			if OneIn(2) {
				obj.Level.SwapActors(obj, other)
				return true, nil
			}
			return true, ErrMoveSwapFailed
		}
	}
	if endtile.Feature == FeatClosedDoor {
		endtile.Feature = FeatOpenDoor
		return true, ErrMoveOpenedDoor
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
		return true, nil
	}

	if obj.IsPlayer() {
		obj.Game.Events.Message(fmt.Sprintf("%v can't go that way.", obj.Spec.Name))
	}

	return conf || false, ErrMoveBlocked
}

// Rest a turn. If p is confused, there is still a chance that it will attempt
// to move. Calling this should always consume a turn.
func (p *ActorMover) Rest() {
	if p.Obj().Sheet.Confused() {
		p.Move(math.Pt(0, 0))
	}
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

// Randomizes a direction.
func confusedir(_ math.Point) math.Point {
	// TODO: Maybe make this less random, and actually dependent on the given
	// point like in Sil's confuse_dir.

	var x, y int

	for {
		x, y = RandInt(-1, 2), RandInt(-1, 2)
		if !(x == 0 && y == 0) {
			break
		}
	}
	return math.Pt(x, y)
}
