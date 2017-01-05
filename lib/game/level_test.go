package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

var lTestActor = &Spec{
	Family:  FamActor,
	Genus:   GenMonster,
	Species: "TestSpecies",
	Name:    "Hi",
	Traits:  &Traits{},
}

var lTestItem = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
	Species: "TestSpecies",
	Name:    "Hiiii",
	Traits:  &Traits{},
}

func TestOkPlaceActor(t *testing.T) {
	g := newTestGame()
	obj := newObj(lTestActor)
	pos := math.Pt(1, 1)

	ok := g.Level.Place(obj, pos)

	if !ok {
		t.Error(`Place((1, 1) was false, want true`)
	}

	if g.Level.At(pos).Actor != obj {
		t.Error(`Place((1, 1)) did not set tile actor to obj`)
	}

	if g.Level.At(pos) != obj.Tile {
		t.Error(`Place((1, 1)) did not set actor's tile to obj`)
	}
}

func TestSecondPlaceActorCleansUp(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestActor)
	startpos := math.Pt(1, 1)
	endpos := math.Pt(1, 2)

	g.Level.Place(obj, startpos)
	g.Level.Place(obj, endpos)

	if g.Level.At(startpos).Actor != nil {
		t.Error(`Place((2, 2)) did not set (1, 1) tile actor to nil`)
	}
	if g.Level.At(endpos).Actor != obj {
		t.Error(`Place((2, 2)) did not set tile actor to obj`)
	}
}

func TestBadActorPlaceOntoSolid(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestActor)
	pos := math.Pt(0, 0)

	ok := g.Level.Place(obj, pos)

	if ok {
		t.Error(`Place( (0,0) ) onto FeatWall ok was true; want false`)
	}
}

func TestBadPlaceActorOntoOccupiedTile(t *testing.T) {
	g := newTestGame()
	a1, a2 := g.NewObj(lTestActor), g.NewObj(lTestActor)
	pos := math.Pt(1, 1)

	g.Level.Place(a1, pos)
	ok := g.Level.Place(a2, pos)

	if ok {
		t.Error(`Place onto other actor: ok was true; want false`)
	}
}

func TestPlaceAddsActorToList(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestActor)
	startpos := math.Pt(1, 1)

	g.Level.Place(obj, startpos)

	if actual := len(g.Level.actors); actual != 2 {
		t.Errorf(`Place(obj) put %d actors in list; want 2`, actual)
	}
	if actual := g.Level.actors[1]; actual != obj {
		t.Errorf(`Place(obj) put %v in list; want %v`, actual, obj)
	}
}

func TestBadPlaceDoesNotAddActorToList(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestActor)
	startpos := math.Pt(0, 0)

	g.Level.Place(obj, startpos)

	// Player is always in the list.
	if actual := len(g.Level.actors); actual != 1 {
		t.Errorf(`Place(obj) put %d actors in list; want 1`, actual)
	}
}

func TestPlaceSingleItem(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestItem)
	dest := math.Pt(1, 1)

	g.Level.Place(obj, dest)

	items := g.Level.At(dest).Items

	if size := items.Len(); size != 1 {
		t.Errorf(`Place(item) put %d items; want 1`, size)
	}

	if item := items.Top(); item != obj {
		t.Errorf(`Place(item) was %v, want %v`, item, obj)
	}
}

func TestOkRemoveActor(t *testing.T) {
	g := newTestGame()
	obj := newObj(lTestActor)
	pos := math.Pt(1, 1)

	g.Level.Place(obj, pos)
	g.Level.Remove(obj)

	if obj.Level != nil {
		t.Error(`obj.Level was not nil`)
	}
	if obj.Tile != nil {
		t.Error(`obj.Tile was not nil`)
	}
	if g.Level.At(pos).Actor != nil {
		t.Error(`Actor's previous tile had tile.Actor != nil`)
	}
	if len(g.Level.actors) > 1 {
		t.Error(`l.actors had monster-actors after removal.`)
	}
}
