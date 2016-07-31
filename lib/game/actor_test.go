package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

type aTestFac struct {
}

var (
	actorTestQueue   = newEventQueue()
	actorTestFactory = &aTestFac{}
)

// Creates a new moving actor for use in these tests.
func (*aTestFac) NewObj(_ *Spec) *Obj {
	spec := &Spec{
		Type:    OTActor,
		Name:    "Hi",
		Subtype: "TestMover",
		Traits:  &Traits{Mover: NewActorMover},
	}
	return newObj(spec, actorTestQueue)
}

func TestOkMove(t *testing.T) {
	l := NewLevel(4, 4, actorTestFactory, IdentLevel)
	obj := actorTestFactory.NewObj(nil)
	startpos := math.Pt(1, 1)

	l.Place(obj, startpos)

	ok := obj.Mover.Move(math.Pt(1, 0))

	if !ok {
		t.Error(`Move( (1, 0)) = false, want true`)
	}

	newpos := obj.Pos()
	want := math.Pt(2, 1)
	if newpos != want {
		t.Errorf(`Move((1, 0)) = %v, want %v`, newpos, want)
	}

	if l.At(startpos).Actor != nil {
		t.Error(`Move((1, 0)) did not set start tile actor to nil`)
	}
	if l.At(newpos).Actor != obj {
		t.Error(`Move((1, 0)) did not set dest tile actor to obj`)
	}
}

func TestActorCollision(t *testing.T) {
	l := NewLevel(4, 4, actorTestFactory, IdentLevel)
	a1, a2 := actorTestFactory.NewObj(nil), actorTestFactory.NewObj(nil)
	l.Place(a1, math.Pt(1, 1))
	l.Place(a2, math.Pt(2, 1))

	ok := a1.Mover.Move(math.Pt(1, 0))

	if ok {
		t.Error(`a1.Move( (1, 0)) = true, want false`)
	}
}

func TestPlayerMaxHPCalc(t *testing.T) {
	obj := &Obj{}
	obj.Stats = &stats{obj: obj, vit: 1}
	obj.Sheet = &PlayerSheet{obj: obj}

	if maxhp, want := obj.Sheet.MaxHP(), 20; maxhp != want {
		t.Errorf(`MaxHP() was %d, want %d`, maxhp, want)
	}
}

func TestPlayerMaxMPCalc(t *testing.T) {
	obj := &Obj{}
	obj.Stats = &stats{obj: obj, mnd: 2}
	obj.Sheet = &PlayerSheet{obj: obj}

	if maxmp, want := obj.Sheet.MaxMP(), 30; maxmp != want {
		t.Errorf(`MaxMP() was %d, want %d`, maxmp, want)
	}
}
