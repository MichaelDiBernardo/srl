package game

import (
	"testing"
)

func TestPositiveMeleeDesc(t *testing.T) {
	g := NewGame()
	spec := &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: "SpecTest",
		Name:    "TestWeapon",
		Traits: &Traits{
			Equip: NewEquip(Equip{
				Melee: 2,
			}),
		},
	}

	item := g.NewObj(spec)

	if desc, want := item.Equip.Desc(), "(+2)"; desc != want {
		t.Errorf("+2 weap Desc() was %v, want %v", desc, want)
	}
}
