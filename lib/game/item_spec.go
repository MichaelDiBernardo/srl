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
		Traits:  &Traits{},
	}
)
