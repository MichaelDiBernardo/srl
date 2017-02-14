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

func restoreEffectsDeps() {
	EffectsSpecs = oldEffectsSpecs
}
