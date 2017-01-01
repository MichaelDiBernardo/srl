package game

// Something you can equip.
const GenEquip = "equip"

// Something you can eat / drink / use a single time.
const GenConsume = "consume"

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

type UseFunc func(user Consumer)

type Consume struct {
	Trait
	Use UseFunc
}

func NewConsume(usefunc UseFunc) func(*Obj) *Consume {
	return func(obj *Obj) *Consume {
		return &Consume{
			Trait: Trait{obj: obj},
			Use:   usefunc,
		}
	}
}
