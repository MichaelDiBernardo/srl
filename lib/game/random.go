package game

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Roll a set of d s-sided die.
func DieRoll(d, s int) int {
	total := 0
	for i := 0; i < d; i++ {
		total += intsource(s + 1)
	}
	return total
}

// Thin wrapper over math/rand so that we can stub out in tests. This is
// basically "what do you want me to use as Intn(n)".
type IntSource func(int) int

// The current source that random functions from here will use.
var intsource IntSource = rand.Intn

// The backup, if the "normal" source has been swapped out.
var oldintsource IntSource = nil

// "Fix" the random generator to return ints from the given sequence.
func FixRandom(ints []int) {
	oldintsource = intsource
	n := 0
	intsource = func(_ int) int {
		x := ints[n]
		n++
		return x
	}
}

// Restore the random generator the way it was.
func RestoreRandom() {
	intsource = oldintsource
	oldintsource = nil
}
