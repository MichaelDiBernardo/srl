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
	fov := make([]math.Point, 0, FOVRadius*FOVRadius)
	cheb := math.Chebyshev(FOVRadius)
	pos, level := a.Obj().Pos(), a.Obj().Level

	for y := cheb.Min.Y; y < cheb.Max.Y; y++ {
		for x := cheb.Min.X; x < cheb.Max.X; x++ {
			pt := math.Pt(x, y).Add(pos)
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
