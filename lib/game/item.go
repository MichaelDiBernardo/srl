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

// See New*Stats in actor.go to understand why this is written this way.
func NewEquip(equip Equip) func(*Obj) *Equip {
	return func(o *Obj) *Equip {
		equip.obj = o
		return &equip
	}
}

func NewNullEquip(o *Obj) *Equip {
	return &Equip{Trait: Trait{obj: o}}
}

func (e *Equip) Desc() string {
	return fmt.Sprintf("(+%v)", e.Melee)
}
