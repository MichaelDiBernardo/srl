package game

import (
	"testing"
)

var stActorSpec = &Spec{
	Family:  FamActor,
	Genus:   GenMonster,
	Species: "TestSpecies",
	Name:    "Hi",
	Traits: &Traits{
		Mover:    NewActorMover,
		Packer:   NewActorPacker,
		Equipper: NewActorEquipper,
		User:     NewActorUser,
		Sheet:    NewPlayerSheet,
	},
}

func TestSkillCheckLost(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)

	FixRandomDie([]int{4, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillCheckWon(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)

	FixRandomDie([]int{3, 1})
	defer RestoreRandom()

	won, by := skillcheck(2, 1, 0, challenger, nil)

	if !won {
		t.Error(`skillcheck(): won is false, want true`)
	}
	if by != 3 {
		t.Errorf(`skillcheck(): by is %d, want 3`, by)
	}
}

func TestSkillCheckTie(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)

	FixRandomDie([]int{1, 1})
	defer RestoreRandom()

	won, by := skillcheck(1, 1, 0, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != 0 {
		t.Errorf(`skillcheck(): by is %d, want 0`, by)
	}
}

func TestSkillCheckResistsAdd10(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)

	FixRandomDie([]int{10, 1, 10, 1, 10, 1})
	defer RestoreRandom()

	won, by := skillcheck(1, 1, 1, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -1 {
		t.Errorf(`skillcheck(): by is %d, want -1`, by)
	}

	won, by = skillcheck(1, 1, 2, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -11 {
		t.Errorf(`skillcheck(): by is %d, want -1`, by)
	}

	won, by = skillcheck(1, 1, -1, challenger, nil)

	if !won {
		t.Error(`skillcheck(): won is false, want true`)
	}
	if by != 19 {
		t.Errorf(`skillcheck(): by is %d, want 19`, by)
	}
}
