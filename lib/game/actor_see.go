package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

var FOVRadius = 4

// A thing that has a FOV
type Seer interface {
	Objgetter
	CalcFOV()
	FOV() FOV
}

type FOV []math.Point

type ActorSeer struct {
	Trait
	fov FOV
}

func NewActorSeer(obj *Obj) Seer {
	return &ActorSeer{Trait: Trait{obj: obj}}
}

func (a *ActorSeer) CalcFOV() {
	// Max area of FOV is (2r+1) * (2r+1)
	maxarea := 4*FOVRadius*FOVRadius + 4*FOVRadius + 1
	fov := make([]math.Point, maxarea)

	pos, level := a.Obj().Pos(), a.Obj().Level

	for r := 0; r <= FOVRadius; r++ {
		edge := math.ChebyEdge(r)
		for _, cpt := range edge {
			pt := pos.Add(cpt)
			if pt.In(level) {
				fov = append(fov, pt)
			}
		}
	}

	a.fov = fov

	if actor := a.Obj(); actor.IsPlayer() {
		actor.Level.UpdateVis()
	}
}

func (a *ActorSeer) FOV() FOV {
	return a.fov
}
