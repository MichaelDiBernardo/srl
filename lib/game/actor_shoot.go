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
	Pos  math.Point
	Path Path
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
		obj.Game.Events.SwitchMode(ModeShoot)
	}
}

func (s *ActorShooter) Targets() []Target {
	return []Target{}
}
