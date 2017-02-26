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

func TestSkillcheckLost(t *testing.T) {
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

func TestSkillcheckWon(t *testing.T) {
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

func TestSkillcheckTie(t *testing.T) {
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

func TestSkillcheckResistsAdd10(t *testing.T) {
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

func TestSkillcheckChallengerCursed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	challenger.Sheet.SetCursed(true)

	// Should pick the 4 over the 5.
	FixRandomDie([]int{4, 5, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillcheckDefenderCursed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	defender := g.NewObj(stActorSpec)
	defender.Sheet.SetCursed(true)

	// Should pick the 7 over the 9.
	FixRandomDie([]int{4, 9, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, defender)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillcheckDefenderBlessedAndCursed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	defender := g.NewObj(stActorSpec)
	defender.Sheet.SetCursed(true)
	defender.Sheet.SetBlessed(true)

	FixRandomDie([]int{4, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, defender)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillcheckChallengerBlessed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	challenger.Sheet.SetBlessed(true)

	// Should pick the 4 over the 2.
	FixRandomDie([]int{4, 2, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, nil)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillcheckDefenderBlessed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	defender := g.NewObj(stActorSpec)
	defender.Sheet.SetBlessed(true)

	// Should pick the 7 over the 5.
	FixRandomDie([]int{4, 5, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, defender)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestSkillcheckDefenderCursedAndBlessed(t *testing.T) {
	g := newTestGame()
	challenger := g.NewObj(stActorSpec)
	defender := g.NewObj(stActorSpec)
	defender.Sheet.SetBlessed(true)
	defender.Sheet.SetCursed(true)

	FixRandomDie([]int{4, 7})
	defer RestoreRandom()

	won, by := skillcheck(1, 2, 0, challenger, defender)

	if won {
		t.Error(`skillcheck(): won is true, want false`)
	}
	if by != -4 {
		t.Errorf(`skillcheck(): by is %d, want -4`, by)
	}
}

func TestCombatrollRollerCursed(t *testing.T) {
	g := newTestGame()
	dude := g.NewObj(stActorSpec)
	dude.Sheet.SetCursed(true)

	// Should pick the 4 over the 9.
	FixRandomDie([]int{9, 4})
	defer RestoreRandom()

	if roll := combatroll(dude); roll != 4 {
		t.Errorf(`combatroll() was %d, want 4`, roll)
	}
}

func TestCombatrollRollerBlessed(t *testing.T) {
	g := newTestGame()
	dude := g.NewObj(stActorSpec)
	dude.Sheet.SetBlessed(true)

	// Should pick the 9 over the 4.
	FixRandomDie([]int{4, 9})
	defer RestoreRandom()

	if roll := combatroll(dude); roll != 9 {
		t.Errorf(`combatroll() was %d, want 9`, roll)
	}
}

func TestCombatrollRollerBlessedAndCursed(t *testing.T) {
	g := newTestGame()
	dude := g.NewObj(stActorSpec)
	dude.Sheet.SetBlessed(true)
	dude.Sheet.SetCursed(true)

	FixRandomDie([]int{9})
	defer RestoreRandom()

	if roll := combatroll(dude); roll != 9 {
		t.Errorf(`combatroll() was %d, want 9`, roll)
	}
}
