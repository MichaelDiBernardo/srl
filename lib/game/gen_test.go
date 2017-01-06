package game

import (
	"testing"
)

func TestNoSpecsGeneratesNoGroups(t *testing.T) {
	groups := Generate(1, 1, 1, []*Spec{}, nil)
	if l := len(groups); l > 0 {
		t.Errorf(`Generate(empty) created %d groups, want 0.`, l)
	}
}

func TestSpecsAtExactDepth(t *testing.T) {
	g := newTestGame()
	specs := []*Spec{
		{
			Name:   "1",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{1},
			},
		},
		{
			Name:   "2",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{2},
			},
		},
		{
			Name:   "3",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{1},
			},
		},
	}

	FixRandomSource([]int{0, 1, 0})
	defer RestoreRandom()

	groups := Generate(3, 1, 0, specs, g)

	if l := len(groups); l != 3 {
		t.Errorf(`Generate() made %d groups, want 3`, l)
	}

	expected := []string{"1", "3", "1"}
	for i, group := range groups {
		if l := len(group); l != 1 {
			t.Errorf(`Generate() made group %d of size %d, want 1`, i, l)
		}
		if s, e := group[0].Spec.Name, expected[i]; s != e {
			t.Errorf(`Wrong spec on obj in group %d; got %s, want %s`, s, e)
		}
	}
}

func TestSpecsWithWiggleDepth(t *testing.T) {
	g := newTestGame()
	specs := []*Spec{
		{
			Name:   "1",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{1},
			},
		},
		{
			Name:   "2",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{2},
			},
		},
		{
			Name:   "3",
			Traits: &Traits{},
			Gen: Gen{
				Depths: []int{1},
			},
		},
	}

	FixRandomSource([]int{0, 1, 0})
	defer RestoreRandom()

	groups := Generate(3, 1, 1, specs, g)

	if l := len(groups); l != 3 {
		t.Errorf(`Generate() made %d groups, want 3`, l)
	}

	expected := []string{"1", "2", "1"}
	for i, group := range groups {
		if l := len(group); l != 1 {
			t.Errorf(`Generate() made group %d of size %d, want 1`, i, l)
		}
		if s, e := group[0].Spec.Name, expected[i]; s != e {
			t.Errorf(`Wrong spec on obj in group %d; got %s, want %s`, s, e)
		}
	}
}

func TestGroupSize(t *testing.T) {
	g := newTestGame()
	specs := []*Spec{
		{
			Name:   "1",
			Traits: &Traits{},
			Gen: Gen{
				Depths:    []int{1},
				GroupSize: 2,
			},
		},
		{
			Name:   "2",
			Traits: &Traits{},
			Gen: Gen{
				Depths:    []int{1},
				GroupSize: 3,
			},
		},
	}

	FixRandomSource([]int{0, 1})
	defer RestoreRandom()

	groups := Generate(2, 1, 0, specs, g)

	if l := len(groups); l != 2 {
		t.Errorf(`Generate() made %d groups, want 3`, l)
	}

	if l := len(groups[0]); l != 2 {
		t.Errorf(`Generate() made group 0 of size %d, want 2`, l)
	}
	if l := len(groups[1]); l != 3 {
		t.Errorf(`Generate() made group 0 of size %d, want 3`, l)
	}
}
