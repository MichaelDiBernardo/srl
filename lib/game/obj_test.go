package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

func TestOkPlace(t *testing.T) {
	obj := NewObj(Traits{})
	l := NewLevel(4, 4, IdentLevel)
	startpos := math.Pt(1, 1)

	ok := obj.Place(l, startpos)

	if !ok {
		t.Error(`Place((1, 1) was false, want true`)
	}

	if l.At(startpos).Actor != obj {
		t.Error(`Place((1, 1)) did not set tile actor to obj`)
	}

	if l.At(startpos) != obj.Tile {
		t.Error(`Place((1, 1)) did not set actor's tile to obj`)
	}
}

func TestSecondPlaceCleansUp(t *testing.T) {
	obj := NewObj(Traits{})
	l := NewLevel(4, 4, IdentLevel)
	startpos := math.Pt(1, 1)
	endpos := math.Pt(2, 2)

	obj.Place(l, startpos)
	obj.Place(l, endpos)

	if l.At(startpos).Actor != nil {
		t.Error(`Place((2, 2)) did not set (1, 1) tile actor to nil`)
	}
	if l.At(endpos).Actor != obj {
		t.Error(`Place((2, 2)) did not set tile actor to obj`)
	}
}

func TestBadPlace(t *testing.T) {
	obj := NewObj(Traits{})
	l := NewLevel(4, 4, SquareLevel)
	startpos := math.Pt(0, 0)

	ok := obj.Place(l, startpos)

	if ok {
		t.Error(`Move( (0,0) ) onto FeatWall ok was true; want false`)
	}
}
