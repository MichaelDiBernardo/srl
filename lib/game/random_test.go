package game

import (
	"testing"
)

func TestRoll(t *testing.T) {
	FixRandom([]int{1, 2, 3, 5})
	defer RestoreRandom()

	if roll, want := DieRoll(4, 1), 11; roll != want {
		t.Errorf(`Die.Roll() was %d, want %d`, roll, want)
	}
}
