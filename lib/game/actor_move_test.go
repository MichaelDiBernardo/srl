package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

func TestOkMove(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(atActorSpec)

	startpos := math.Pt(1, 1)
	g.Level.Place(obj, startpos)

	ok := obj.Mover.Move(math.Pt(1, 0))

	if !ok {
		t.Error(`Move( (1, 0)) = false, want true`)
	}

	newpos := obj.Pos()
	want := math.Pt(2, 1)
	if newpos != want {
		t.Errorf(`Move((1, 0)) = %v, want %v`, newpos, want)
	}

	if g.Level.At(startpos).Actor != nil {
		t.Error(`Move((1, 0)) did not set start tile actor to nil`)
	}
	if g.Level.At(newpos).Actor != obj {
		t.Error(`Move((1, 0)) did not set dest tile actor to obj`)
	}
}

func TestMonsterSwapping(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	g.Level.Place(a1, math.Pt(1, 1))
	g.Level.Place(a2, math.Pt(2, 1))

	ok := a1.Mover.Move(math.Pt(1, 0))

	if !ok {
		t.Error(`a1.Move( (1, 0)) = false, want true`)
	}

	if a1.Pos() != math.Pt(2, 1) && a2.Pos() != math.Pt(1, 1) {
		t.Error(`Monsters did not swap`)
	}
}

func TestMoveOpensClosedDoor(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(atActorSpec)

	startpos := math.Pt(1, 1)
	doorpos := math.Pt(1, 2)
	g.Level.Place(obj, startpos)
	g.Level.At(doorpos).Feature = FeatClosedDoor

	ok := obj.Mover.Move(math.Pt(0, 1))

	if ok {
		t.Error(`Move into closed door was ok = true, want false`)
	}

	if feat := g.Level.At(doorpos).Feature; feat != FeatOpenDoor {
		t.Errorf(`Door didn't open; got feature %#v, want %#v`, feat, FeatOpenDoor)
	}
}
