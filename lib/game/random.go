package game

import (
	"fmt"
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
		total += 1 + intsource(s)
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

// "Fix" the random generator to return ints from the given sequence. Also
// returns the closure that substitutes intsource, if you want to chain it in
// subsequent fakes.
func FixRandomSource(ints []int) func(int) int {
	oldintsource = intsource
	n := 0
	intsource = func(_ int) int {
		x := ints[n]
		n++
		return x
	}
	return intsource
}

// An intsource that subtracts 1 to each element in the list. This makes it
// compatible with rigging dierolls directly, since DieRoll has to add one to
// each int to represent a roll from 1 to n (instead of 0 to n-1).
func FixRandomDie(ints []int) func(int) int {
	f := FixRandomSource(ints)
	intsource = func(i int) int {
		return f(i) - 1
	}
	return intsource
}

// Restore the random generator the way it was.
func RestoreRandom() {
	intsource = oldintsource
	oldintsource = nil
}

type Dice struct {
	Dice  int
	Sides int
}

var ZeroDice = Dice{}

func NewDice(dice, sides int) Dice {
	return Dice{Dice: dice, Sides: sides}
}

func (d Dice) Add(dice, sides int) Dice {
	return NewDice(d.Dice+dice, d.Sides+sides)
}

func (d Dice) Roll() int {
	return DieRoll(d.Dice, d.Sides)
}

func (d Dice) String() string {
	return fmt.Sprintf("%dd%d", d.Dice, d.Sides)
}

func RandInt(low int, high int) int {
	return low + intsource(high-low)
}

func OneIn(n int) bool {
	return intsource(n) == 0
}

func Coinflip() bool {
	return OneIn(2)
}

type Weighter interface {
	Weight() int
}

// Selects 1 item from a weighted list of choices. In the weighted list {a: 1,
// b: 2, c: 1}, we'd expect to see a selected 25% of the time, b selected 50%
// of the time, and c selected 25% of the time.
func WChoose(choices []Weighter) (pos int, chosen Weighter) {
	if len(choices) == 0 {
		return -1, nil
	}

	// This is not necessary -- the loop below would handle this case. However,
	// we use this function REALLY often in fighting-related tests where we
	// already have to stub out a bunch of dierolls, and basically all of those
	// tests use a single attack. To avoid having to put a leading '0' in every
	// fake rolls list, we handle this case without consuming a random number.
	if len(choices) == 1 {
		return 0, choices[0]
	}

	// Compute total weight of all items in choices.
	tw := 0
	for _, item := range choices {
		tw += item.Weight()
	}

	// Select a random number in [0,tw). Now, continually add up weights. When
	// we finally get a weight > r, select that item.
	r, cw := RandInt(0, tw), 0
	for i, item := range choices {
		cw += item.Weight()
		if r < cw {
			return i, item
		}
	}

	// This should never happen.
	panic(fmt.Sprintf("Could not WChoose from %+v", choices))
}
