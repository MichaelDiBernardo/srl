package game

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Thin wrapper over math/rand so that we can stub out in tests. This is
// basically "what do you want me to use as Intn(n)".
type Random func(int) int

// Represents a set of nds die.
type Die struct {
	Dice  int
	Sides int
	r     Random
}

// Create new die.
func NewDie(dice, sides int) Die {
	return Die{
		Dice:  dice,
		Sides: sides,
		r:     rand.Intn,
	}
}

func (d Die) Roll() int {
	roll := 0
	for i := 0; i < d.Dice; i++ {
		roll += d.r(d.Sides + 1)
	}
	return roll
}
