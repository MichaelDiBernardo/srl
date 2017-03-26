package game

import (
	"testing"
)

func TestSeenMonsterGivesNoXP(t *testing.T) {
	g := newTestGame()
	mon := g.NewObj(atActorSpec)
	mon.Seen = true

	g.Player.Learner.GainXPSight(mon)

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

	g.Player.Learner.GainXPSight(mon)

	// Should get 30 xp for first sighting.
	if xp := g.Player.Learner.XP(); xp != 30 {
		t.Errorf(`Learner.XP() was %d, want 30`, xp)
	}

	g.Player.Learner.GainXPSight(mon)

	// Should get an additional 15 for next.
	if xp := g.Player.Learner.XP(); xp != 45 {
		t.Errorf(`Learner.XP() was %d, want 45`, xp)
	}

	g.Player.Learner.GainXPSight(mon)

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

	g.Player.Learner.GainXPSight(item)

	if xp := g.Player.Learner.XP(); xp != 20 {
		t.Errorf(`Learner.XP() was %d, want 20`, xp)
	}

	g.Player.Learner.GainXPSight(item)

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

	g.Player.Learner.GainXPKill(mon)

	// Should get 30 xp for first sighting.
	if xp := g.Player.Learner.XP(); xp != 30 {
		t.Errorf(`Learner.XP() was %d, want 30`, xp)
	}

	g.Player.Learner.GainXPKill(mon)

	// Should get an additional 15 for next.
	if xp := g.Player.Learner.XP(); xp != 45 {
		t.Errorf(`Learner.XP() was %d, want 45`, xp)
	}

	g.Player.Learner.GainXPKill(mon)

	// Should get an additional 10 for next.
	if xp := g.Player.Learner.XP(); xp != 55 {
		t.Errorf(`Learner.XP() was %d, want 52`, xp)
	}
}

func TestBeginLearningCalledTwiceError(t *testing.T) {
	g := newTestGame()

	_, err := g.Player.Learner.BeginLearning()

	if err != nil {
		t.Errorf(`First call to BeginLearning() returned error %v; want nil`, err)
	}

	_, err = g.Player.Learner.BeginLearning()

	if err != ErrAlreadyLearning {
		t.Errorf(`Second call to BeginLearning() returned %v; want %v`, err, ErrAlreadyLearning)
	}
}

func TestCancelLearningBeforeBeginError(t *testing.T) {
	g := newTestGame()

	err := g.Player.Learner.CancelLearning()

	if err != ErrNotLearning {
		t.Errorf(`CancelLearning() returned error %v; want %v`, err, ErrNotLearning)
	}
}

func TestEndLearningBeforeBeginError(t *testing.T) {
	g := newTestGame()

	err := g.Player.Learner.EndLearning()

	if err != ErrNotLearning {
		t.Errorf(`EndLearning() returned error %v; want %v`, err, ErrNotLearning)
	}
}

func TestLearningScenario(t *testing.T) {
	g := newTestGame()
	l, s := g.Player.Learner, g.Player.Sheet

	// Set up some base skill investments and mods so that we can test that the
	// cost only relies on the unmodded skill amount.
	s.SetSkill(Melee, 1)
	s.SetSkill(Evasion, 2)
	s.SetSkill(Shooting, 3)
	s.ChangeSkillMod(Melee, 1)
	s.ChangeSkillMod(Evasion, 2)
	s.ChangeSkillMod(Shooting, 3)

	// Give the player some initial XP to spend.
	l.(*ActorLearner).gainxp(5000)

	// Learn some stuff.
	l.BeginLearning()
	l.LearnSkill(Melee)            // 200
	l.LearnSkill(Evasion)          // 300
	c, _ := l.LearnSkill(Shooting) // 400

	// Player changes
	if sk, want := s.UnmodSkill(Melee), 2; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if sk, want := s.UnmodSkill(Evasion), 3; sk != want {
		t.Errorf(`s.UnmodSkill(Evasion) was %d, want %d`, sk, want)
	}
	if sk, want := s.UnmodSkill(Shooting), 4; sk != want {
		t.Errorf(`s.UnmodSkill(Shooting) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 4100; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	// Change record upkeep.
	if cost, want := c.TotalCost, 900; cost != want {
		t.Errorf(`c.TotalCost was %d, want %d`, cost, want)
	}
	if change, want := c.Changes[Melee], (SkillChangeItem{Cost: 200, Points: 1}); change != want {
		t.Errorf(`c.Changes[Melee] was %+v, want %+v`, change, want)
	}
	if change, want := c.Changes[Evasion], (SkillChangeItem{Cost: 300, Points: 1}); change != want {
		t.Errorf(`c.Changes[Evasion] was %+v, want %+v`, change, want)
	}
	if change, want := c.Changes[Shooting], (SkillChangeItem{Cost: 400, Points: 1}); change != want {
		t.Errorf(`c.Changes[Shooting] was %+v, want %+v`, change, want)
	}

	c, _ = l.LearnSkill(Melee)
	if sk, want := s.UnmodSkill(Melee), 3; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 3800; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	if cost, want := c.TotalCost, 1200; cost != want {
		t.Errorf(`c.TotalCost was %d, want %d`, cost, want)
	}
	if change, want := c.Changes[Melee], (SkillChangeItem{Cost: 500, Points: 2}); change != want {
		t.Errorf(`c.Changes[Melee] was %+v, want %+v`, change, want)
	}

	l.UnlearnSkill(Melee)
	c, _ = l.UnlearnSkill(Melee)

	if sk, want := s.UnmodSkill(Melee), 1; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 4300; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	if cost, want := c.TotalCost, 700; cost != want {
		t.Errorf(`c.TotalCost was %d, want %d`, cost, want)
	}
	if change, want := c.Changes[Melee], (SkillChangeItem{}); change != want {
		t.Errorf(`c.Changes[Melee] was %+v, want %+v`, change, want)
	}

	l.EndLearning()

	// Check player's final state after learning.
	if sk, want := s.UnmodSkill(Melee), 1; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if sk, want := s.UnmodSkill(Evasion), 3; sk != want {
		t.Errorf(`s.UnmodSkill(Evasion) was %d, want %d`, sk, want)
	}
	if sk, want := s.UnmodSkill(Shooting), 4; sk != want {
		t.Errorf(`s.UnmodSkill(Shooting) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 4300; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	if totalxp, want := l.TotalXP(), 5000; totalxp != want {
		t.Errorf(`l.TotalXP() was %d, want %d`, totalxp, want)
	}
}

func TestCancelLearning(t *testing.T) {
	g := newTestGame()
	l, s := g.Player.Learner, g.Player.Sheet
	l.(*ActorLearner).gainxp(5000)

	l.BeginLearning()
	l.LearnSkill(Melee)
	l.LearnSkill(Melee)
	l.LearnSkill(Melee)

	l.CancelLearning()

	if sk, want := s.UnmodSkill(Melee), 0; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 5000; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
}

func TestCantOverspendXP(t *testing.T) {
	g := newTestGame()
	l, s := g.Player.Learner, g.Player.Sheet
	l.(*ActorLearner).gainxp(299)

	l.BeginLearning()
	l.LearnSkill(Melee)
	c, err := l.LearnSkill(Melee)

	if want := ErrNotEnoughXP; err != want {
		t.Errorf(`l.LearnSkill() returned err %v, want %v`, err, want)
	}
	if sk, want := s.UnmodSkill(Melee), 1; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 199; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	if cost, want := c.TotalCost, 100; cost != want {
		t.Errorf(`c.TotalCost was %d, want %d`, cost, want)
	}
	if change, want := c.Changes[Melee], (SkillChangeItem{Cost: 100, Points: 1}); change != want {
		t.Errorf(`c.Changes[Melee] was %+v, want %+v`, change, want)
	}
}

func TestCantOverRefundXP(t *testing.T) {
	g := newTestGame()
	l, s := g.Player.Learner, g.Player.Sheet
	l.(*ActorLearner).gainxp(100)

	l.BeginLearning()
	l.LearnSkill(Melee)
	l.UnlearnSkill(Melee)
	c, err := l.UnlearnSkill(Melee)

	if want := ErrNoPointsLearned; err != want {
		t.Errorf(`l.LearnSkill() returned err %v, want %v`, err, want)
	}
	if sk, want := s.UnmodSkill(Melee), 0; sk != want {
		t.Errorf(`s.UnmodSkill(Melee) was %d, want %d`, sk, want)
	}
	if xp, want := l.XP(), 100; xp != want {
		t.Errorf(`l.XP() was %d, want %d`, xp, want)
	}
	if cost, want := c.TotalCost, 0; cost != want {
		t.Errorf(`c.TotalCost was %d, want %d`, cost, want)
	}
	if change, want := c.Changes[Melee], (SkillChangeItem{}); change != want {
		t.Errorf(`c.Changes[Melee] was %+v, want %+v`, change, want)
	}
}
