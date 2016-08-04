package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

type lTestFac struct {
}

var lTestQueue = newEventQueue()

var lTestActor = &Spec{
	Family:  FamActor,
	Genus:   GenMonster,
	Species: "TestSpecies",
	Name:    "Hi",
	Traits:  &Traits{},
}

var lTestItem = &Spec{
	Family:  FamItem,
	Genus:   GenEquip,
	Species: "TestSpecies",
	Name:    "Hiiii",
	Traits:  &Traits{},
}

func TestOkPlaceActor(t *testing.T) {
	l := NewLevel(4, 4, nil, IdentLevel)
	obj := newObj(lTestActor, lTestQueue)
	pos := math.Pt(1, 1)

	ok := l.Place(obj, pos)

	if !ok {
		t.Error(`Place((1, 1) was false, want true`)
	}

	if l.At(pos).Actor != obj {
		t.Error(`Place((1, 1)) did not set tile actor to obj`)
	}

	if l.At(pos) != obj.Tile {
		t.Error(`Place((1, 1)) did not set actor's tile to obj`)
	}
}

func TestSecondPlaceActorCleansUp(t *testing.T) {
	l := NewLevel(4, 4, nil, IdentLevel)
	obj := newObj(lTestActor, lTestQueue)
	startpos := math.Pt(1, 1)
	endpos := math.Pt(2, 2)

	l.Place(obj, startpos)
	l.Place(obj, endpos)

	if l.At(startpos).Actor != nil {
		t.Error(`Place((2, 2)) did not set (1, 1) tile actor to nil`)
	}
	if l.At(endpos).Actor != obj {
		t.Error(`Place((2, 2)) did not set tile actor to obj`)
	}
}

func TestBadActorPlaceOntoSolid(t *testing.T) {
	l := NewLevel(4, 4, nil, SquareLevel)
	obj := newObj(lTestActor, lTestQueue)
	pos := math.Pt(0, 0)

	ok := l.Place(obj, pos)

	if ok {
		t.Error(`Place( (0,0) ) onto FeatWall ok was true; want false`)
	}
}

func TestBadPlaceActorOntoOccupiedTile(t *testing.T) {
	l := NewLevel(4, 4, nil, IdentLevel)
	a1, a2 := newObj(lTestActor, lTestQueue), newObj(lTestActor, lTestQueue)
	pos := math.Pt(0, 0)

	l.Place(a1, pos)
	ok := l.Place(a2, pos)

	if ok {
		t.Error(`Place onto other actor: ok was true; want false`)
	}
}

func TestPlaceAddsActorToList(t *testing.T) {
	l := NewLevel(4, 4, nil, IdentLevel)
	obj := newObj(lTestActor, lTestQueue)
	startpos := math.Pt(1, 1)

	l.Place(obj, startpos)

	if actual := len(l.actors); actual != 1 {
		t.Errorf(`Place(obj) put %d actors in list; want 1`, actual)
	}
	if actual := l.actors[0]; actual != obj {
		t.Errorf(`Place(obj) put %v in list; want %v`, actual, obj)
	}
}

func TestBadPlaceDoesNotAddActorToList(t *testing.T) {
	l := NewLevel(4, 4, nil, SquareLevel)
	obj := newObj(lTestActor, lTestQueue)
	startpos := math.Pt(0, 0)

	l.Place(obj, startpos)

	if actual := len(l.actors); actual != 0 {
		t.Errorf(`Place(obj) put %d actors in list; want 0`, actual)
	}
}

func TestPlaceSingleItem(t *testing.T) {
	l := NewLevel(4, 4, nil, SquareLevel)
	obj := newObj(lTestItem, lTestQueue)
	dest := math.Pt(1, 1)

	l.Place(obj, dest)

	items := l.At(dest).Items

	if size := items.Len(); size != 1 {
		t.Errorf(`Place(item) put %d items; want 1`)
	}

	if item := items.Top(); item != obj {
		t.Errorf(`Place(item) was %v, want %v`, item, obj)
	}
}
