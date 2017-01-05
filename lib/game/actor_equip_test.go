package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

func TestTryEquipWithNoEquipsInInventory(t *testing.T) {
	g := newTestGame()
	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()
	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryEquip w no equips switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryEquipWithEquipsInInventory(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	equipper.Packer.Inventory().Add(equip)

	equipper.Equipper.TryEquip()

	if mode := g.mode; mode != ModeEquip {
		t.Errorf(`TryEquip switched to mode %v, want %v`, mode, ModeEquip)
	}
}

func TestEquipIntoEmptySlot(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	inv := equipper.Packer.Inventory()
	inv.Add(equip)

	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Was mode %v after equip; want %v`, mode, ModeHud)
	}

	if !inv.Empty() {
		t.Errorf(`Item did not leave inventory after equipping.`)
	}

	slot := equip.Equipment.Slot
	if equipped := equipper.Equipper.Body().Slots[slot]; equipped != equip {
		t.Errorf(`Equipped item was %v, want %v`, equipped, equip)
	}
}

func TestEquipIntoOccupiedSlot(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip1 := g.NewObj(atItemSpec)
	equip2 := g.NewObj(atItemSpec)

	inv := equipper.Packer.Inventory()
	inv.Add(equip1)
	inv.Add(equip2)

	// Wield equip1
	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)
	// Wield equip2, swapping out equip1
	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)

	if swapped := inv.Top(); swapped != equip1 {
		t.Errorf(`First wield was not swapped out; got %v, want %v.`, swapped, equip1)
	}

	slot := equip2.Equipment.Slot
	if equipped := equipper.Equipper.Body().Slots[slot]; equipped != equip2 {
		t.Errorf(`Equipped item was %v, want %v`, equipped, equip2)
	}
}

func TestEquipOutOfBounds(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	inv := equipper.Packer.Inventory()
	inv.Add(equip)

	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Was mode %v after equip; want %v`, mode, ModeHud)
	}
}

func TestTryRemoveNothingEquipped(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeHud)
	}
}

func TestTryRemoveSomethingEquipped(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeRemove {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeRemove)
	}
}

func TestRemove(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(equip.Equipment.Slot)

	if removed := equipper.Equipper.Body().Slots[equip.Equipment.Slot]; removed != nil {
		t.Errorf(`Found %v in removed slot; want nil`, removed)
	}

	if removed := equipper.Packer.Inventory().Top(); removed != equip {
		t.Errorf(`Found %v in pack; want %v`, removed, equip)
	}
}

func TestRemoveOutOfBounds(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(100)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Out-of-bounds removed switched mode to %v; want %v`, mode, ModeHud)
	}
}

func TestRemoveOverflowsToGround(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Packer.Inventory().capacity = 0
	g.Level.Place(equipper, math.Pt(1, 1))
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(equip.Equipment.Slot)

	if removed := equipper.Packer.Inventory().Top(); removed != nil {
		t.Errorf(`Found %v in pack; want nil`, removed)
	}

	if removed := equipper.Tile.Items.Top(); removed != equip {
		t.Errorf(`Found %v on floor; want %v`, removed, equip)
	}
}

func TestRemoveWhenInvAndGroundAreFull(t *testing.T) {
	g := newTestGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)
	equipper.Equipper.Body().Wear(equip)

	g.Level.Place(equipper, math.Pt(1, 1))

	equipper.Packer.Inventory().capacity = 0
	equipper.Tile.Items.capacity = 0

	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeHud)
	}

}
