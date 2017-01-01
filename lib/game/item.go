package game

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
type ConsumeFunc func(consumer Consumer)

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
