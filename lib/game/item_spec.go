package game

const (
	SpecFist         = "fist"
	SpecSword        = "sword"
	SpecLeatherArmor = "leatherarmor"

	SpecCure = "cure"
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
				Effects: Effects{BrandFire},
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
				Effects:  Effects{ResistPoison},
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
}
