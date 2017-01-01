package game

// Something you can equip.
const GenEquip = "equip"

// Something you can eat / drink / use a single time.
const GenConsume = "consume"

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

type ConsumeFunc func(consumer Consumer)

// Consumable trait.
type Consumable struct {
	Trait
	Consume ConsumeFunc
}

func NewConsumable(cf ConsumeFunc) func(*Obj) *Consumable {
	return func(obj *Obj) *Consumable {
		return &Consumable{
			Trait:   Trait{obj: obj},
			Consume: cf,
		}
	}
}
