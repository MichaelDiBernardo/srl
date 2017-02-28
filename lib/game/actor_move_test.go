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

	if _, err := obj.Mover.Move(math.Pt(1, 0)); err != nil {
		t.Errorf(`Move( (1, 0)) = %v, want true`, err)
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

func TestMonsterSwapWorks(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	g.Level.Place(a1, math.Pt(1, 1))
	g.Level.Place(a2, math.Pt(2, 1))

	// There's a random failure rate attached to swapping. This forces swap.
	FixRandomSource([]int{0})
	defer RestoreRandom()

	if _, err := a1.Mover.Move(math.Pt(1, 0)); err != nil {
		t.Errorf(`a1.Move( (1, 0)) = %v, want nil`, err)
	}

	if a1.Pos() != math.Pt(2, 1) && a2.Pos() != math.Pt(1, 1) {
		t.Error(`Monsters did not swap`)
	}
}

func TestMonsterSwapFails(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	g.Level.Place(a1, math.Pt(1, 1))
	g.Level.Place(a2, math.Pt(2, 1))

	// There's a random failure rate attached to swapping. This forces swap to
	// fail.
	FixRandomSource([]int{1})
	defer RestoreRandom()

	if _, err := a1.Mover.Move(math.Pt(1, 0)); err != ErrMoveSwapFailed {
		t.Errorf(`a1.Move( (1, 0)) = %v, want %v`, err, ErrMoveSwapFailed)
	}

	if a1.Pos() != math.Pt(1, 1) && a2.Pos() != math.Pt(2, 1) {
		t.Error(`Monsters swapped!`)
	}
}

func TestCantSwapWithPetrifiedTarget(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	a2.Sheet.SetPetrified(true)

	g.Level.Place(a1, math.Pt(1, 1))
	g.Level.Place(a2, math.Pt(2, 1))

	if _, err := a1.Mover.Move(math.Pt(1, 0)); err != ErrMoveSwapFailed {
		t.Errorf(`a1.Move( (1, 0)) = %v, want %v`, err, ErrMoveSwapFailed)
	}

	if a1.Pos() != math.Pt(1, 1) && a2.Pos() != math.Pt(2, 1) {
		t.Error(`Monsters swapped!`)
	}
}

func TestMoveOpensClosedDoor(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(atActorSpec)

	startpos := math.Pt(1, 1)
	doorpos := math.Pt(1, 2)
	g.Level.Place(obj, startpos)
	g.Level.At(doorpos).Feature = FeatClosedDoor

	if _, err := obj.Mover.Move(math.Pt(0, 1)); err != ErrMoveOpenedDoor {
		t.Errorf(`Move into closed door was %v, want %v`, err, ErrMoveOpenedDoor)
	}

	if feat := g.Level.At(doorpos).Feature; feat != FeatOpenDoor {
		t.Errorf(`Door didn't open; got feature %#v, want %#v`, feat, FeatOpenDoor)
	}
}
