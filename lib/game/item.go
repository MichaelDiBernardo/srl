package game

import (
	"fmt"
)

const GenEquip = "equip"

type Equip struct {
	Trait
	Damroll  Dice
	Protroll Dice
	Melee    int
	Evasion  int
	Weight   int
	Slot     Slot
}

// See NewSheet in actor.go to understand why this is written this way.
func NewEquip(equipspec Equip) func(*Obj) *Equip {
	return func(o *Obj) *Equip {
		// Copy.
		equip := equipspec
		equip.obj = o
		return &equip
	}
}

func (e *Equip) Desc() string {
	return fmt.Sprintf("(+%v)", e.Melee)
}
