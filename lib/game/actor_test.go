package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

func TestPlayerOkMove(t *testing.T) {
	obj := NewObj(Traits{Mover: NewActorMover})
	l := NewLevel(4, 4, IdentLevel)
	startpos := math.Pt(1, 1)
	obj.Place(l, startpos)

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
