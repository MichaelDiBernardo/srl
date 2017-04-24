package game

import (
	"errors"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

var ErrTargetOutOfRange = errors.New("TargetOutOfRange")
var ErrNoClearShot = errors.New("NoClearShot")

// Anything that can attack at range.
type Shooter interface {
	Objgetter
	// Primarily for the player -- try to switch to "shooting mode". If the
	// player has no ranged weapon or is otherwise in no shape to shoot, this
	// will emit a message and cancel the modeswitch.
	TryShoot()
	// Return a list of targets in LOS, sorted by proximity.
	Targets() []Target
	// Given a point on the map, this will give the target info for shooting at
	// that spot. Will return ErrTargetOutOfRange or ErrNoClearShot if the shot
	// is impossible.
	Target(math.Point) (Target, error)
}

// A target in LOS of the shooter. Points in Pos and Path are relative to the
// actor, who is considered to be at origin. So, if there is a target to the
// east with one space intervening, its Pos would be (2,0) and its path would
// be {(1,0), (2,0))
type Target struct {
	Pos    math.Point
	Path   Path
	Target *Obj
}

type ActorShooter struct {
	Trait
}

func NewActorShooter(obj *Obj) Shooter {
	return &ActorShooter{Trait: Trait{obj: obj}}
}

func (s *ActorShooter) TryShoot() {
	obj := s.obj
	if obj.Equipper.Body().Shooter() == nil {
		obj.Game.Events.Message("Nothing to shoot with.")
	} else if s.obj.Sheet.Afraid() {
		msg := fmt.Sprintf("%s is too afraid to shoot!", obj.Spec.Name)
		obj.Game.Events.Message(msg)
	} else if s.obj.Sheet.Confused() {
		msg := fmt.Sprintf("%s is too confused to shoot!", obj.Spec.Name)
		obj.Game.Events.Message(msg)
	} else {
		obj.Game.SwitchMode(ModeShoot)
	}
}

func (s *ActorShooter) Targets() []Target {
	fov := s.obj.Senser.FOV()
	targets := []Target{}

	for _, p := range fov {
		tile := s.obj.Game.Level.At(p)
		victim := tile.Actor

		if victim == nil || victim == s.obj {
			continue
		}

		target, err := s.Target(tile.Pos)
		if err != nil {
			continue
		}
		targets = append(targets, target)
	}
	return targets
}

func (s *ActorShooter) Target(p math.Point) (Target, error) {
	// TODO: Replace with sheet attack range, or something.
	mypos, lev := s.obj.Pos(), s.obj.Game.Level
	srange := s.obj.Equipper.Body().Shooter().Equipment.Range

	if math.EucDist(s.obj.Pos(), p) > srange {
		return Target{}, ErrTargetOutOfRange
	}

	path, ok := lev.FindPath(mypos, p, PathCost)
	if !ok {
		return Target{}, ErrNoClearShot
	}

	target := Target{
		Pos:    p,
		Path:   path,
		Target: lev.At(p).Actor,
	}
	return target, nil
}
