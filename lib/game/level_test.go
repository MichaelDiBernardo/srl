package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

var lTestActor = &Spec{
	Family:  FamActor,
	Genus:   GenMonster,
	Species: "TestSpecies",
	Gen:     Gen{Floors: []int{2}},
	Name:    "Hi",
	Traits:  &Traits{Sheet: NewPlayerSheet},
}

var lTestItem = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
	Species: "TestSpecies2",
	Gen:     Gen{Floors: []int{1}},
	Name:    "Hiiii",
	Traits:  &Traits{},
}

var lTestItem2 = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
	Species: "TestSpecies3",
	Gen:     Gen{Floors: []int{3}},
	Name:    "Hiiiiiiii",
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

	if actual := g.Level.scheduler.Len(); actual != 2 {
		t.Errorf(`Place(obj) put %d actors in list; want 2`, actual)
	}
}

func TestBadPlaceDoesNotAddActorToList(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(lTestActor)
	startpos := math.Pt(0, 0)

	g.Level.Place(obj, startpos)

	// Player is always in the list.
	if actual := g.Level.scheduler.Len(); actual != 1 {
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
	if g.Level.scheduler.Len() > 1 {
		t.Error(`l.actors had monster-actors after removal.`)
	}
}

type pftest struct {
	pic   string
	start math.Point
	end   math.Point
	ok    bool
	want  Path
}

func TestPathfinding(t *testing.T) {
	tests := []pftest{
		{
			start: math.Pt(0, 0),
			end:   math.Pt(1, 0),
			ok:    true,
			want:  Path{math.Pt(1, 0)},
			pic: `
@x`,
		},
		{
			start: math.Pt(0, 0),
			end:   math.Pt(0, 0),
			ok:    true,
			want:  Path{},
		},
		{
			start: math.Pt(-1, -1),
			end:   math.Pt(0, 0),
			ok:    false,
			want:  Path{math.Pt(1, 0)},
		},
		{
			start: math.Pt(0, 0),
			end:   math.Pt(6, 0),
			ok:    true,
			want:  Path{math.Pt(1, 0), math.Pt(2, 0), math.Pt(3, 0), math.Pt(4, 0), math.Pt(5, 0), math.Pt(6, 0)},
			pic: `
@     x
#######`,
		},
		{
			start: math.Pt(0, 0),
			end:   math.Pt(3, 3),
			ok:    true,
			want:  Path{math.Pt(1, 0), math.Pt(2, 0), math.Pt(3, 1), math.Pt(3, 2), math.Pt(3, 3)},
			pic: `
@   # 
### #
### #
###x#
`,
		},
		{
			start: math.Pt(0, 1),
			end:   math.Pt(2, 1),
			ok:    true,
			want:  Path{math.Pt(1, 0), math.Pt(2, 1)},
			pic: `
#'# 
@#x
#+#
`,
		},
		{
			start: math.Pt(0, 1),
			end:   math.Pt(2, 1),
			ok:    false,
			want:  Path{},
			pic: `
###
@#x
###`,
		},
	}

	for ti, test := range tests {
		g := newTestGame()
		l := NewLevel(40, 40, g, StringLevel(test.pic))
		g.Level = l

		path, ok := l.FindPath(test.start, test.end, PathCost)

		if ok != test.ok {
			t.Errorf(`Pathfinding test %d: ok=%v, want=%v`, ti, ok, test.ok)
		}
		if !test.ok {
			continue
		}

		if alen, wlen := len(path), len(test.want); alen != wlen {
			t.Errorf(`Pathfinding test %d: len(actual)=%d, len(want)=%d`, ti, alen, wlen)
			continue
		}

		for i, wpt := range test.want {
			if apt := path[i]; apt != wpt {
				t.Errorf(`Pathfinding test %d: mismatch at %d, got %v, want %v`, ti, i, path, test.want)
				break
			}
		}
	}
}

func TestUpdateVisTeachesPlayer(t *testing.T) {
	g := newTestGame()
	dest := math.Pt(1, 1)

	i1, i2, m := g.NewObj(lTestItem), g.NewObj(lTestItem2), g.NewObj(lTestActor)

	g.Level.Place(g.Player, dest)
	g.Level.Place(i1, dest)
	g.Level.Place(i2, dest)

	mloc := math.Pt(1, 2)
	g.Level.Place(m, mloc)

	g.Level.UpdateVis(Field{dest, mloc})

	if xp := g.Player.Learner.XP(); xp != 60 {
		t.Errorf(`Learner.XP() was %d, want 60`, xp)
	}

	for i, obj := range []*Obj{i1, i2, m} {
		if !obj.Seen {
			t.Errorf(`o.Seen (%d, %v) was false, want true`, i, obj)
		}
	}
}
