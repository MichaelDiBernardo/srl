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
