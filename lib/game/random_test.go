package game

import (
	"testing"
)

func TestRoll(t *testing.T) {
	FixRandomDie([]int{1, 2, 3, 5})
	defer RestoreRandom()

	if roll, want := DieRoll(4, 1), 11; roll != want {
		t.Errorf(`Die.Roll() was %d, want %d`, roll, want)
	}
}

type wchooseT struct {
	w int
}

func (t wchooseT) Weight() int {
	return t.w
}

type wchooseTC struct {
	items []Weighter
	ints  []int
	pos   int
}

func TestWChoose(t *testing.T) {
	tests := []wchooseTC{
		wchooseTC{
			items: []Weighter{wchooseT{10}},
			ints:  []int{},
			pos:   0,
		},
		wchooseTC{
			items: []Weighter{wchooseT{1}, wchooseT{1}},
			ints:  []int{0},
			pos:   0,
		},
		wchooseTC{
			items: []Weighter{wchooseT{1}, wchooseT{1}},
			ints:  []int{1},
			pos:   1,
		},
		wchooseTC{
			items: []Weighter{wchooseT{2}, wchooseT{1}},
			ints:  []int{1},
			pos:   0,
		},
		wchooseTC{
			items: []Weighter{wchooseT{2}, wchooseT{10}, wchooseT{1}},
			ints:  []int{12},
			pos:   2,
		},
		wchooseTC{
			items: []Weighter{},
			ints:  []int{},
			pos:   -1,
		},
	}

	for ti, test := range tests {
		func() {
			FixRandomSource(test.ints)
			defer RestoreRandom()

			pos, _ := WChoose(test.items)

			if pos != test.pos {
				t.Errorf(`Test %d: WChoose() selected %d, want %d`, ti, pos, test.pos)
			}
		}()
	}
}
