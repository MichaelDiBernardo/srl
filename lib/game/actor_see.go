package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

var LOSRadius = 4

// A thing that has a LOS
type Seer interface {
	Objgetter
	CalcLOS()
	LOS() LOS
}

type LOS []math.Point

type ActorSeer struct {
	Trait
	los LOS
}

func NewActorSeer(obj *Obj) Seer {
	return &ActorSeer{Trait: Trait{obj: obj}}
}

func (a *ActorSeer) CalcLOS() {
	a.los = make([]math.Point, 0)
}

func (a *ActorSeer) LOS() LOS {
	return a.los
}
