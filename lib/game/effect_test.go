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
	fakeVuln1
)

var (
	oldResistMap = ResistMap
	oldVulnMap   = VulnMap
	oldBrands    = Brands
)

func TestHas(t *testing.T) {
	effects := Effects{fakeEffect1}

	if !effects.Has(fakeEffect1) {
		t.Error(`effects.Has(fakeEffect1) was false, want true`)
	}
	if effects.Has(fakeEffect2) {
		t.Error(`effects.Has(fakeEffect2) was true, want false`)
	}
}

func TestResists(t *testing.T) {
	effects := Effects{fakeResist1}
	ResistMap = map[Effect]Effect{fakeEffect1: fakeResist1}
	defer restoreEffectsDeps()

	if !effects.Resists(fakeEffect1) {
		t.Error(`effects.Resists(fakeEffect1) was false, want true`)
	}
	if effects.Resists(fakeEffect2) {
		t.Error(`effects.Resists(fakeEffect1) was true, want false`)
	}
}

func TestVulnTo(t *testing.T) {
	effects := Effects{fakeVuln1}
	VulnMap = map[Effect]Effect{fakeEffect1: fakeVuln1}
	defer restoreEffectsDeps()

	if !effects.VulnTo(fakeEffect1) {
		t.Error(`effects.VulnTo(fakeEffect1) was false, want true`)
	}
	if effects.VulnTo(fakeEffect2) {
		t.Error(`effects.Resists(fakeEffect2) was true, want false`)
	}
}

func TestBrands(t *testing.T) {
	effects := Effects{fakeBrand1, fakeEffect1, fakeBrand2, fakeResist1}
	Brands = Effects{fakeBrand1, fakeBrand2}
	defer restoreEffectsDeps()

	brands := effects.Brands()
	if l := len(brands); l != 2 {
		t.Errorf(`len(effects.Brands()) was %d, want 2`, l)
	}

	for i, b := range brands {
		if b != Brands[i] {
			t.Errorf(`effects.Brands(): Didn't find brand %d at %d, got %d`, b, i, Brands[i])
		}
	}
}

func restoreEffectsDeps() {
	ResistMap = oldResistMap
	VulnMap = oldVulnMap
	Brands = oldBrands
}
