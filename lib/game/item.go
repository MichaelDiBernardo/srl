package game

import (
	"fmt"
)

// Something you can equip.
const GenEquipment = "equipment"

// Something you can eat / drink / use a single time.
const GenConsumable = "consume"

// Equipment trait.
type Equipment struct {
	Trait
	Damroll  Dice
	Protroll Dice
	Melee    int
	Evasion  int
	Weight   int
	Slot     Slot
	Effects  Effects
}

// See NewSheet in actor.go to understand why this is written this way.
func NewEquipment(equipspec Equipment) func(*Obj) *Equipment {
	return func(o *Obj) *Equipment {
		// Copy.
		equip := equipspec
		equip.obj = o
		return &equip
	}
}

// Function that actually does something when this item gets used.
type ConsumeFunc func(user User)

// Consumable trait.
type Consumable struct {
	Trait
	Consume ConsumeFunc
}

// Given a consumefunc, this creates a factory function for consumables with
// this action.
func NewConsumable(cf ConsumeFunc) func(*Obj) *Consumable {
	return func(obj *Obj) *Consumable {
		return &Consumable{
			Trait:   Trait{obj: obj},
			Consume: cf,
		}
	}
}

func curefunc(user User) {
	u := user.Obj()
	u.Sheet.Heal(40)
	u.Game.Events.Message(fmt.Sprintf("%s recovers.", u.Spec.Name))
}

func stimfunc(user User) {
	u := user.Obj()
	u.Game.Events.Message(fmt.Sprintf("%s is wracked with pain.", u.Spec.Name))
	u.Sheet.Hurt(DieRoll(4, 4))
	u.Ticker.AddEffect(EffectStim, DieRoll(20, 4))
}
