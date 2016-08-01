package game

import (
	"testing"
)

func TestRoll(t *testing.T) {
	rolls, n := [4]int{1, 2, 3, 5}, 0
	r := func(_ int) int {
		x := rolls[n]
		n++
		return x
	}

	sut := NewDie(4, 1)
	sut.r = r

	if roll, want := sut.Roll(), 11; roll != want {
		t.Errorf(`Die.Roll() was %d, want %d`, roll, want)
	}
}
