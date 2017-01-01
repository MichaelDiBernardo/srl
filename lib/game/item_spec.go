package game

const (
	SpecFist         = "fist"
	SpecSword        = "sword"
	SpecLeatherArmor = "leatherarmor"

	SpecCure = "cure"
)

var (
	WeapSword = &Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecSword,
		Name:    "SWORD",
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Damroll: NewDice(2, 5),
				Melee:   0,
				Evasion: 1,
				Weight:  2,
				Slot:    SlotHand,
			}),
		},
	}
	ArmorLeather = &Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
		Species: SpecLeatherArmor,
		Name:    "LEATHER",
		Traits: &Traits{
			Equipment: NewEquipment(Equipment{
				Protroll: NewDice(1, 4),
				Melee:    0,
				Evasion:  -1,
				Weight:   4,
				Slot:     SlotBody,
			}),
		},
	}
)

var (
	PotCure = &Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecCure,
		Name:    "CURE",
		Traits: &Traits{
			Consumable: NewConsumable(curefunc),
		},
	}
)
