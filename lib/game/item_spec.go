package game

const (
	SpecSword = "sword"
)

var (
	WeapSword = &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: SpecSword,
		Name:    "SWORD",
		Traits: &Traits{
			Equip: NewEquip(Equip{
				Damroll: NewDice(2, 5),
				Melee:   0,
				Evasion: 1,
				Weight:  2,
				Slot:    SlotHand,
			}),
		},
	}
)
