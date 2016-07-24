package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
	"math/rand"
)

type ActorSpec struct {
	Type   ObjSubtype
	Traits *Traits
}

func NewActor(spec *ActorSpec) *Obj {
	obj := NewObj(spec.Traits)
	obj.Type = OTActor
	obj.Subtype = spec.Type
	return obj
}

// A thing that can move given a specific direction.
type Mover interface {
	Move(dir math.Point) bool
}

type ActorMover struct {
	obj *Obj
}

func NewActorMover(obj *Obj) Mover {
	return &ActorMover{obj: obj}
}

// Try to move the player. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	return obj.Level.Place(obj, endpos)
}

// A thing that can move given a specific direction.
type AI interface {
	Act(l *Level) bool
}

type RandomAI struct {
	obj *Obj
}

func NewRandomAI(obj *Obj) AI {
	return &RandomAI{obj: obj}
}

func (ai *RandomAI) Act(l *Level) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	dir := math.Pt(x, y)
	log.Printf("AI: Moving from %v by %v", ai.obj.Pos(), dir)

	return ai.obj.Mover.Move(dir)
}
