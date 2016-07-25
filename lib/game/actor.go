package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
	"math/rand"
)

// A specification for a type of actor.
type ActorSpec struct {
	Type   ObjSubtype
    Name   string
	Traits *Traits
}

// Create an actor from its specification. If you're making these for use in an
// actual game, you should use game.NewActor(...) instead.
func newActor(spec *ActorSpec, eq *EventQueue) *Obj {
	return newObj(OTActor, spec.Type, spec.Traits, eq)
}

// A thing that can move given a specific direction.
type Mover interface {
	Move(dir math.Point) bool
}

// A dummy mover used in cases where a thing can't move.
type nullMover struct {
}

// Do nothing and return false.
func (_ *nullMover) Move(dir math.Point) bool {
	return false
}

// Singleton instance of the null mover.
var theNullMover = &nullMover{}

// Constructor for null movers.
func NewNullMover(_ *Obj) Mover {
	return theNullMover
}

// A universally-applicable mover for actors.
type ActorMover struct {
	obj *Obj
}

// Constructor for actor movers.
func NewActorMover(obj *Obj) Mover {
	return &ActorMover{obj: obj}
}

// Try to move the actor. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

    moved := obj.Level.Place(obj, endpos)
    if !moved {
        p.obj.Events.Message("You can't go that way!")
    }
    return moved
}

// A thing that can move given a specific direction.
type AI interface {
	Act(l *Level) bool
}

// A dummy AI used when a thing doesn't need a computer to think for it.
type nullAI struct {
}

// Do nothing and return false.
func (_ *nullAI) Act(l *Level) bool {
	return false
}

// Singleton instance of the null mover.
var theNullAI = &nullAI{}

// Constructor for null movers.
func NewNullAI(_ *Obj) AI {
	return theNullAI
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	obj *Obj
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{obj: obj}
}

// Move in any of the 8 directions with uniform chance. Does not take walls
// etc. in account so this will happily try to bump into things.
func (ai *RandomAI) Act(l *Level) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	dir := math.Pt(x, y)
	log.Printf("AI: Moving from %v by %v", ai.obj.Pos(), dir)

	return ai.obj.Mover.Move(dir)
}
