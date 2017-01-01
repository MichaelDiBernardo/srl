package game

const (
	SpecFist         = "fist"
	SpecSword        = "sword"
	SpecLeatherArmor = "leatherarmor"
)

var (
	WeapSword = &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
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
		Genus:   GenEquip,
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
