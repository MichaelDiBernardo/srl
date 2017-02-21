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
