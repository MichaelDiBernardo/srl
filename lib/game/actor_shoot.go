package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Anything that can attack at range.
type Shooter interface {
	Objgetter
	// Primarily for the player -- try to switch to "shooting mode". If the
	// player has no ranged weapon or is otherwise in no shape to shoot, this
	// will emit a message and cancel the modeswitch.
	TryShoot()
	// Return a list of targets in LOS, sorted by proximity.
	Targets() []Target
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
	fov, lev, pos := s.obj.Senser.FOV(), s.obj.Game.Level, s.obj.Pos()
	targets := []Target{}

	for _, p := range fov {
		tile := lev.At(p)
		victim := tile.Actor

		if victim == nil || victim == s.obj {
			continue
		}

		path, ok := lev.FindPath(pos, tile.Pos, PathCost)
		if !ok {
			continue
		}
		for i := 0; i < len(path); i++ {
			path[i] = path[i].Sub(pos)
		}

		target := Target{
			Pos:    tile.Pos.Sub(pos),
			Path:   path,
			Target: victim,
		}
		targets = append(targets, target)
	}
	return targets
}
