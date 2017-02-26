package game

import (
	"testing"
)

const (
	fakeEffect1 Effect = NumEffects + iota
	fakeEffect2
	fakeBrand1
	fakeBrand2
	fakeResist1
	fakeResist2
)

var oldEffectsSpecs = EffectsSpecs
var testEffectsSpecs = EffectsSpec{
	fakeEffect1: {Type: EffectTypeStatus, ResistedBy: fakeResist1},
	fakeEffect2: {Type: EffectTypeStatus, ResistedBy: fakeResist2},
	fakeBrand1:  {Type: EffectTypeBrand},
	fakeBrand2:  {Type: EffectTypeBrand},
	fakeResist1: {Type: EffectTypeResist},
	fakeResist2: {Type: EffectTypeResist},
}

func TestHas(t *testing.T) {
	EffectsSpecs = testEffectsSpecs
	defer restoreEffectsDeps()

	effects := NewEffects(map[Effect]int{fakeEffect1: 1})

	if n := effects.Has(fakeEffect1); n != 1 {
		t.Errorf(`effects.Has(fakeEffect1) was %d, want 1`, n)
	}
	if n := effects.Has(fakeEffect2); n != 0 {
		t.Errorf(`effects.Has(fakeEffect2) was %d, want 0`, n)
	}
}

func TestResists(t *testing.T) {
	EffectsSpecs = testEffectsSpecs
	defer restoreEffectsDeps()

	effects := NewEffects(map[Effect]int{fakeResist1: 1, fakeResist2: -1})

	if n := effects.Resists(fakeEffect1); n != 1 {
		t.Errorf(`effects.Resists(fakeEffect1) was %d, want 1`, n)
	}
	if n := effects.Resists(fakeEffect2); n != -1 {
		t.Errorf(`effects.Resists(fakeEffect1) was %d, want -1`, n)
	}
	if n := effects.Resists(fakeBrand1); n != 0 {
		t.Errorf(`effects.Resists(fakeEffect1) was %d, want 0`, n)
	}
}

func TestBrands(t *testing.T) {
	EffectsSpecs = testEffectsSpecs
	defer restoreEffectsDeps()

	effects := NewEffects(map[Effect]int{
		fakeBrand1:  1,
		fakeEffect1: 1,
		fakeBrand2:  1,
		fakeResist1: 2,
	})

	brands := effects.Brands()
	if l := len(brands); l != 2 {
		t.Errorf(`len(effects.Brands()) was %d, want 2`, l)
	}

	for e, info := range brands {
		if info.Type != EffectTypeBrand {
			t.Errorf(`Effect %v had type %v, want %v`, e, info.Type, EffectTypeBrand)
		}
	}
}

func TestMerge(t *testing.T) {
	EffectsSpecs = testEffectsSpecs
	defer restoreEffectsDeps()

	e1 := NewEffects(map[Effect]int{fakeBrand1: 1})
	e2 := NewEffects(map[Effect]int{fakeEffect1: 1})

	merged := e1.Merge(e2)
	// We should not modify e1 or e2 to produce merged.
	if l := len(merged); l != 2 {
		t.Errorf(`len(e1.Merge(e2)) was %d, want %d`, l, 2)
	}
	if l := len(e1); l != 1 {
		t.Errorf(`e1 changed after merge: new length %d`, l)
	}
	if l := len(e2); l != 1 {
		t.Errorf(`e2 changed after merge: new length %d`, l)
	}

	// Are the elements of merged ok?
	if f1, f2 := merged[fakeBrand1], e1[fakeBrand1]; f1 != f2 {
		t.Errorf(`merged[fakeBrand1] != e1[fakeBrand1]; %v vs %v`, f1, f2)
	}
	if f1, f2 := merged[fakeEffect1], e2[fakeEffect1]; f1 != f2 {
		t.Errorf(`merged[fakeEffect1] != e2[fakeEffect1]; %v vs %v`, f1, f2)
	}

	// Merge merged with e1 to see if count increments.
	merged = merged.Merge(e1)
	if l := len(merged); l != 2 {
		t.Errorf(`len(e1.Merge(e2)) was %d, want %d`, l, 2)
	}
	if newcount := merged[fakeBrand1].Count; newcount != 2 {
		t.Errorf(`merged[fakeBrand1].Count was %d, want 2`, newcount)
	}
	if count := merged[fakeEffect1].Count; count != 1 {
		t.Errorf(`merged[fakeEffect1].Count was %d, want 1`, count)
	}
}

func TestResistDmg(t *testing.T) {
	EffectsSpecs = testEffectsSpecs
	defer restoreEffectsDeps()

	effects := NewEffects(map[Effect]int{fakeResist1: 1, fakeResist2: -1})

	if n := effects.ResistDmg(fakeEffect1, 10); n != 5 {
		t.Errorf(`effects.ResistDmg(fakeEffect1) was %d, want 5`, n)
	}
	if n := effects.ResistDmg(fakeEffect2, 10); n != 20 {
		t.Errorf(`effects.ResistDmg(fakeEffect1) was %d, want 20`, n)
	}

	effects = NewEffects(map[Effect]int{fakeResist1: 2, fakeResist2: -3})

	if n := effects.ResistDmg(fakeEffect1, 10); n != 3 {
		t.Errorf(`effects.ResistDmg(fakeEffect1) was %d, want 3`, n)
	}
	if n := effects.ResistDmg(fakeEffect2, 10); n != 40 {
		t.Errorf(`effects.ResistDmg(fakeEffect1) was %d, want 30`, n)
	}
}

func restoreEffectsDeps() {
	EffectsSpecs = oldEffectsSpecs
}

func TestRegen0(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 1}},
	})
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
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 1}},
		regen: 1,
	})
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
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 1}},
		regen: 2,
	})
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
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 1}},
		hp:    0,
		regen: 1,
	})
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

func TestActiveStun(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectStun, 10)
	if lvl := obj.Sheet.Stun(); lvl != Stunned {
		t.Errorf(`Stunlevel was %v, want %v`, lvl, Stunned)
	}
	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != -2 {
			t.Errorf(`Skill %v had mod %d after stun, want -2`, skill, s)
		}
	}

	obj.Ticker.AddEffect(EffectStun, 40)
	if lvl := obj.Sheet.Stun(); lvl != MoreStunned {
		t.Errorf(`Stunlevel was %v, want %v`, lvl, MoreStunned)
	}
	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != -4 {
			t.Errorf(`Skill %v had mod %d after stun, want -4`, skill, s)
		}
	}

	// The value of 'delay' doesn't matter at all here, what matters is that a
	// turn has passed. This should reduce the stun penalty to -2.
	obj.Ticker.Tick(0)
	if lvl := obj.Sheet.Stun(); lvl != Stunned {
		t.Errorf(`Stunlevel was %v, want %v`, lvl, Stunned)
	}
	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != -2 {
			t.Errorf(`Skill %v had mod %d after stun, want -2`, skill, s)
		}
	}

	for i := 0; i < 49; i++ {
		obj.Ticker.Tick(0)
	}
	if lvl := obj.Sheet.Stun(); lvl != NotStunned {
		t.Errorf(`Stunlevel was %v, want %v`, lvl, NotStunned)
	}
	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != 0 {
			t.Errorf(`Skill %v had mod %d after stun, want 0`, skill, s)
		}
	}
}

func TestActivePoison(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})
	hpstart := obj.Sheet.HP()

	obj.Ticker.AddEffect(EffectPoison, 10)

	obj.Ticker.Tick(0)
	if lost := hpstart - obj.Sheet.HP(); lost != 2 {
		t.Errorf(`Ticking poison counter inflicted %d, want %d`, lost, 2)
	}

	obj.Ticker.Tick(0)
	if lost := hpstart - obj.Sheet.HP(); lost != 3 {
		t.Errorf(`Ticking poison counter twice inflicted %d, want %d`, lost, 3)
	}

	for i := 0; i < 7; i++ {
		obj.Ticker.Tick(0)
	}

	if lost := hpstart - obj.Sheet.HP(); lost != 10 {
		t.Errorf(`Ticking poison counter to end inflicted %d, want %d`, lost, 10)
	}
	if left := obj.Ticker.Counter(EffectPoison); left != 0 {
		t.Errorf(`Poison counter at %d, want 0`, left)
	}
}

func TestActiveCut(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})
	hpstart := obj.Sheet.HP()

	obj.Ticker.AddEffect(EffectCut, 10)

	obj.Ticker.Tick(0)
	if lost := hpstart - obj.Sheet.HP(); lost != 2 {
		t.Errorf(`Ticking cut counter inflicted %d, want %d`, lost, 2)
	}

	obj.Ticker.Tick(0)
	if lost := hpstart - obj.Sheet.HP(); lost != 3 {
		t.Errorf(`Ticking cut counter twice inflicted %d, want %d`, lost, 3)
	}

	for i := 0; i < 7; i++ {
		obj.Ticker.Tick(0)
	}

	if lost := hpstart - obj.Sheet.HP(); lost != 10 {
		t.Errorf(`Ticking cut counter to end inflicted %d, want %d`, lost, 10)
	}
	if left := obj.Ticker.Counter(EffectCut); left != 0 {
		t.Errorf(`Poison counter at %d, want 0`, left)
	}
}

func TestActiveBlind(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectBlind, 1)

	if !obj.Sheet.Blind() {
		t.Error(`obj.Sheet.Blind() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Blind() {
		t.Error(`obj.Sheet.Blind() was true, want false`)
	}
}

func TestActiveSlow(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectSlow, 1)

	if !obj.Sheet.Slow() {
		t.Error(`obj.Sheet.Slow() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Slow() {
		t.Error(`obj.Sheet.Slow() was true, want false`)
	}
}

func TestActiveConfuse(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectConfuse, 1)

	if !obj.Sheet.Confused() {
		t.Error(`obj.Sheet.Confused() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Confused() {
		t.Error(`obj.Sheet.Confused() was true, want false`)
	}
}

func TestActiveFear(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectFear, 1)

	if !obj.Sheet.Afraid() {
		t.Error(`obj.Sheet.Afraid() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Afraid() {
		t.Error(`obj.Sheet.Afraid() was true, want false`)
	}
}

func TestActivePara(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectPara, 1)

	if !obj.Sheet.Paralyzed() {
		t.Error(`obj.Sheet.Paralyzed() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Paralyzed() {
		t.Error(`obj.Sheet.Paralyzed() was true, want false`)
	}
}

func TestActiveSilence(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectSilence, 1)

	if !obj.Sheet.Silenced() {
		t.Error(`obj.Sheet.Silenced() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Silenced() {
		t.Error(`obj.Sheet.Silenced() was true, want false`)
	}
}

func TestActiveCursed(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	obj.Ticker.AddEffect(EffectCurse, 1)

	if !obj.Sheet.Cursed() {
		t.Error(`obj.Sheet.Cursed() was false, want true`)
	}

	obj.Ticker.Tick(0)

	if obj.Sheet.Cursed() {
		t.Error(`obj.Sheet.Cursed() was true, want false`)
	}
}

func TestActiveStim(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	// Stimming twice should still only give one buff.
	obj.Ticker.AddEffect(EffectStim, 1)
	obj.Ticker.AddEffect(EffectStim, 1)

	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != 2 {
			t.Errorf(`Skill %v had mod %d after stim, want 2`, skill, s)
		}
	}

	obj.Ticker.Tick(0)
	obj.Ticker.Tick(0)
	for skill := Melee; skill < NumSkills; skill++ {
		if s := obj.Sheet.SkillMod(skill); s != 0 {
			t.Errorf(`Skill %v had mod %d after stun, want 0`, skill, s)
		}
	}
}

func TestActiveHyper(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	// Stimming twice should still only give one buff.
	obj.Ticker.AddEffect(EffectHyper, 1)
	obj.Ticker.AddEffect(EffectHyper, 1)

	for stat := Str; stat < NumStats; stat++ {
		if s := obj.Sheet.StatMod(stat); s != 2 {
			t.Errorf(`Stat %v had mod %d after stat, want 2`, stat, s)
		}
	}

	obj.Ticker.Tick(0)
	obj.Ticker.Tick(0)
	for stat := Str; stat < NumStats; stat++ {
		if s := obj.Sheet.StatMod(stat); s != 0 {
			t.Errorf(`Stat %v had mod %d after stat, want 0`, stat, s)
		}
	}
}
