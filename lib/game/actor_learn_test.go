package game

import (
	"testing"
)

func TestSeenMonsterGivesNoXP(t *testing.T) {
	g := newTestGame()
	mon := g.NewObj(atActorSpec)
	mon.Seen = true

	g.Player.Learner.LearnSight(mon)

	if xp := g.Player.Learner.XP(); xp != 0 {
		t.Errorf(`Learner.XP() was %d, want 0`, xp)
	}
}

func TestSeenXPDecaysForMonsters(t *testing.T) {
	monspec := &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Gen: Gen{
			Floors: []int{3},
		},
		Traits: &Traits{},
	}
	g := newTestGame()
	mon := g.NewObj(monspec)

	g.Player.Learner.LearnSight(mon)

	// Should get 30 xp for first sighting.
	if xp := g.Player.Learner.XP(); xp != 30 {
		t.Errorf(`Learner.XP() was %d, want 30`, xp)
	}

	g.Player.Learner.LearnSight(mon)

	// Should get an additional 15 for next.
	if xp := g.Player.Learner.XP(); xp != 45 {
		t.Errorf(`Learner.XP() was %d, want 45`, xp)
	}

	g.Player.Learner.LearnSight(mon)

	// Should get an additional 10 for next.
	if xp := g.Player.Learner.XP(); xp != 55 {
		t.Errorf(`Learner.XP() was %d, want 52`, xp)
	}
}

func TestSeenItemXPIs0AfterFirstSighting(t *testing.T) {
	itemspec := &Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: "testspec",
		Gen: Gen{
			Floors: []int{2},
		},
		Traits: &Traits{},
	}
	g := newTestGame()
	item := g.NewObj(itemspec)

	g.Player.Learner.LearnSight(item)

	if xp := g.Player.Learner.XP(); xp != 20 {
		t.Errorf(`Learner.XP() was %d, want 20`, xp)
	}

	g.Player.Learner.LearnSight(item)

	if xp := g.Player.Learner.XP(); xp != 20 {
		t.Errorf(`Learner.XP() was %d, want 20`, xp)
	}
}

func TestKillXPDecaysForMonsters(t *testing.T) {
	monspec := &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Gen: Gen{
			Floors: []int{3},
		},
		Traits: &Traits{},
	}
	g := newTestGame()
	mon := g.NewObj(monspec)

	g.Player.Learner.LearnKill(mon)

	// Should get 30 xp for first sighting.
	if xp := g.Player.Learner.XP(); xp != 30 {
		t.Errorf(`Learner.XP() was %d, want 30`, xp)
	}

	g.Player.Learner.LearnKill(mon)

	// Should get an additional 15 for next.
	if xp := g.Player.Learner.XP(); xp != 45 {
		t.Errorf(`Learner.XP() was %d, want 45`, xp)
	}

	g.Player.Learner.LearnKill(mon)

	// Should get an additional 10 for next.
	if xp := g.Player.Learner.XP(); xp != 55 {
		t.Errorf(`Learner.XP() was %d, want 52`, xp)
	}
}
