package game

import (
	"testing"
)

func TestRegen0(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, vit: 1, hp: 0, regen: 0}
	// The regen period should heal 0 HP.
	delay := GetDelay(2) * RegenPeriod
	obj.Ticker.Tick(delay)

	if hp := obj.Sheet.HP(); hp != 0 {
		t.Errorf(`Regen healed %d, want 0`, hp)
	}
}

func TestRegen1(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, vit: 1, hp: 0, regen: 1}
	// Half the regen period should heal 50% HP
	delay := GetDelay(2) * RegenPeriod / 2
	obj.Ticker.Tick(delay)

	if hp := obj.Sheet.HP(); hp != 10 {
		t.Errorf(`Regen healed %d, want 10`, hp)
	}
}

func TestRegen2(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, vit: 1, hp: 0, regen: 2}
	// Quarter the regen period should heal 50% HP
	delay := GetDelay(2) * RegenPeriod / 4
	obj.Ticker.Tick(delay)

	if hp := obj.Sheet.HP(); hp != 10 {
		t.Errorf(`Regen healed %d, want 10`, hp)
	}
}

func TestRegenAcrossLevels(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, vit: 1, hp: 0, regen: 1}
	// If we regen a long time on one floor, and then a shorter time on the
	// next, it should be the same as if we'd regened everything on the same
	// floor. (We're talking about floors here because the total delay counter
	// resets on every floor. This is simulating what happens when you take
	// your first turn on the next floor, but we use large delays to make the
	// intent of the math easier to understand.
	obj.Ticker.Tick(GetDelay(2) * RegenPeriod / 2)
	obj.Ticker.Tick(GetDelay(2) * RegenPeriod / 4)

	if hp := obj.Sheet.HP(); hp != 15 {
		t.Errorf(`Regen healed %d, want 15`, hp)
	}
}
