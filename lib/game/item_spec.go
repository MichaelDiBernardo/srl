package game

const (
	SpecFist         = "fist"
	SpecSword        = "sword"
	SpecColt         = "colt"
	SpecLeatherArmor = "leatherarmor"

	SpecCure    = "cure"
	SpecStim    = "stim"
	SpecHyper   = "hyper"
	SpecRestore = "restore"
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
				Hit:     0,
				Evasion: 1,
				Weight:  3,
				Slot:    SlotHand,
				Effects: NewEffects(map[Effect]int{EffectVamp: 1, SlayPearl: 1}),
			}),
		},
	},
	&Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecColt,
		Name:    "COLT",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Damroll: NewDice(6, 1),
				Hit:     0,
				Evasion: 0,
				Weight:  3,
				Slot:    SlotRanged,
				Range:   5,
				Effects: NewEffects(map[Effect]int{}),
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
				Hit:      0,
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
	&Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecHyper,
		Name:    "HYPER",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Consumable: NewConsumable(hyperfunc),
		},
	},
	&Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecRestore,
		Name:    "RESTORE",
		Gen: Gen{
			Floors:    []int{1},
			GroupSize: 1,
		},
		Traits: &Traits{
			Consumable: NewConsumable(restorefunc),
		},
	},
}
