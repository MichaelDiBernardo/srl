package game

const (
	SpecFist         = "fist"
	SpecSword        = "sword"
	SpecLeatherArmor = "leatherarmor"

	SpecCure = "cure"
	SpecStim = "stim"
)

var Items = []*Spec{
	&Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecSword,
		Name:    "SWORD",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Damroll: NewDice(2, 5),
				Melee:   0,
				Evasion: 1,
				Weight:  2,
				Slot:    SlotHand,
				Effects: NewEffects(map[Effect]int{EffectPara: 1}),
			}),
		},
	},
	&Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecLeatherArmor,
		Name:    "LEATHER",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Protroll: NewDice(1, 4),
				Melee:    0,
				Evasion:  -1,
				Weight:   4,
				Slot:     SlotBody,
				Effects:  NewEffects(map[Effect]int{ResistFire: 1, ResistPoison: 1, ResistStun: 1}),
			}),
		},
	},
	&Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecCure,
		Name:    "CURE",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Consumable: NewConsumable(curefunc),
		},
	},
	&Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecStim,
		Name:    "STIM",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Consumable: NewConsumable(stimfunc),
		},
	},
}
