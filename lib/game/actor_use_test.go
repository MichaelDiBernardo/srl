package game

import (
	"testing"
)

var atConsumeSpec = &Spec{
	Family:  FamItem,
	Genus:   GenConsumable,
	Species: SpecCure,
	Name:    "CURE",
	Traits: &Traits{
		Consumable: NewConsumable(func(u User) {}),
	},
}

func TestTryUseWithNoUsablesInInventory(t *testing.T) {
	g := newTestGame()
	user := g.NewObj(atActorSpec)
	user.User.TryUse()
	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryUse w no usables switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryUseWithConsumablesInInventory(t *testing.T) {
	g := newTestGame()
	user := g.NewObj(atActorSpec)

	item := g.NewObj(atConsumeSpec)
	user.Packer.Inventory().Add(item)
	user.User.TryUse()

	if mode := g.mode; mode != ModeUse {
		t.Errorf(`TryUse w usables switched to mode %v, want %v`, mode, ModeUse)
	}
}

func TestUse(t *testing.T) {
	g := newTestGame()
	user := g.NewObj(atActorSpec)
	used := false

	cspec := &Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecCure,
		Name:    "CURE",
		Traits: &Traits{
			Consumable: NewConsumable(func(u User) {
				used = true
			}),
		},
	}

	item := g.NewObj(cspec)
	user.Packer.Inventory().Add(item)
	user.User.TryUse()
	user.User.Use(0)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Using item switched to mode %v, want %v`, mode, ModeHud)
	}
	if !used {
		t.Error(`Use() did not use item.`)
	}
	if !user.Packer.Inventory().Empty() {
		t.Error(`Use() did not consume item.`)
	}
}

func TestUseOutOfBounds(t *testing.T) {
	g := newTestGame()
	user := g.NewObj(atActorSpec)
	item := g.NewObj(atConsumeSpec)

	user.Packer.Inventory().Add(item)
	user.User.TryUse()
	user.User.Use(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Using item switched to mode %v, want %v`, mode, ModeHud)
	}
	if user.Packer.Inventory().Empty() {
		t.Error(`Use() out-of-bounds consumed first item.`)
	}
}

func TestUseBadItem(t *testing.T) {
	g := newTestGame()
	user := g.NewObj(atActorSpec)

	cspec := &Spec{
		Family:  FamItem,
		Genus:   GenConsumable,
		Species: SpecCure,
		Name:    "CURE",
		Traits:  &Traits{},
	}

	// A potion.
	item := g.NewObj(cspec)
	// A weapon.
	equip := g.NewObj(atItemSpec)
	user.Packer.Inventory().Add(item)
	user.Packer.Inventory().Add(equip)
	user.User.TryUse()
	user.User.Use(1)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Using item switched to mode %v, want %v`, mode, ModeHud)
	}
	if nitems := user.Packer.Inventory().Len(); nitems != 2 {
		t.Error(`Use() consumed nonitem.`)
	}
}
