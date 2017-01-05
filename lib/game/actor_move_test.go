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

func TestActorCollision(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	g.Level.Place(a1, math.Pt(1, 1))
	g.Level.Place(a2, math.Pt(2, 1))

	ok := a1.Mover.Move(math.Pt(1, 0))

	if ok {
		t.Error(`a1.Move( (1, 0)) = true, want false`)
	}
}
