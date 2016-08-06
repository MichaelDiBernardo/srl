package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

var actorTestSpec = &Spec{
	Family:  FamActor,
	Genus:   GenMonster,
	Species: "TestSpecies",
	Name:    "Hi",
	Traits:  &Traits{Mover: NewActorMover},
}

func TestOkMove(t *testing.T) {
	g := NewGame()
	l := NewLevel(4, 4, g, IdentLevel)
	obj := g.NewObj(actorTestSpec)
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
	g := NewGame()
	l := NewLevel(4, 4, g, IdentLevel)
	a1, a2 := g.NewObj(actorTestSpec), g.NewObj(actorTestSpec)
	l.Place(a1, math.Pt(1, 1))
	l.Place(a2, math.Pt(2, 1))

	ok := a1.Mover.Move(math.Pt(1, 0))

	if ok {
		t.Error(`a1.Move( (1, 0)) = true, want false`)
	}
}

func TestPlayerMaxHPCalc(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Stats = &stats{Trait: Trait{obj: obj}, vit: 1}
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}}

	if maxhp, want := obj.Sheet.MaxHP(), 20; maxhp != want {
		t.Errorf(`MaxHP() was %d, want %d`, maxhp, want)
	}
}

func TestPlayerMaxMPCalc(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Stats = &stats{Trait: Trait{obj: obj}, mnd: 2}
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}}

	if maxmp, want := obj.Sheet.MaxMP(), 30; maxmp != want {
		t.Errorf(`MaxMP() was %d, want %d`, maxmp, want)
	}
}

type fakefighter struct {
	Trait
	Called bool
}

func (f *fakefighter) Hit(other Fighter) {
	f.Called = true
}

func (_ *fakefighter) MeleeRoll() int {
	return 0
}

func (_ *fakefighter) EvasionRoll() int {
	return 0
}

func (_ *fakefighter) DamRoll() int {
	return 0
}

func (_ *fakefighter) ProtRoll() int {
	return 0
}

func TestPlayerMonsterCollisionsHit(t *testing.T) {
	g := NewGame()
	player := g.NewObj(PlayerSpec)
	pf := &fakefighter{Trait: Trait{obj: player}}
	player.Fighter = pf

	monster := g.NewObj(actorTestSpec)
	mf := &fakefighter{Trait: Trait{obj: player}}
	monster.Fighter = mf

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(player, math.Pt(0, 0))
	l.Place(monster, math.Pt(1, 1))

	player.Mover.Move(math.Pt(1, 1))

	if !pf.Called {
		t.Error("Moving player into other did not try to hit.")
	}

	monster.Mover.Move(math.Pt(-1, -1))

	if !mf.Called {
		t.Error("Moving other into player did not try to hit.")
	}
}

func TestMonsterMonsterCollisionsHit(t *testing.T) {
	g := NewGame()
	mon1 := g.NewObj(actorTestSpec)
	mf1 := &fakefighter{Trait: Trait{obj: mon1}}
	mon1.Fighter = mf1

	mon2 := g.NewObj(actorTestSpec)
	mf2 := &fakefighter{Trait: Trait{obj: mon2}}
	mon2.Fighter = mf2

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(mon1, math.Pt(0, 0))
	l.Place(mon2, math.Pt(1, 1))

	mon1.Mover.Move(math.Pt(1, 1))

	if mf1.Called {
		t.Error("Moving monster into monster tried to hit.")
	}
}
