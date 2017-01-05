package game

import (
	"testing"
)

var invTestEquip = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
	Species: "testspec",
	Name:    "Item",
	Traits:  &Traits{},
}

var invTestConsumable = &Spec{
	Family:  FamItem,
	Genus:   GenConsumable,
	Species: "testspec",
	Name:    "Item",
	Traits:  &Traits{},
}

func TestTake(t *testing.T) {
	g := newTestGame()
	inv := NewInventory()

	item1 := g.NewObj(invTestEquip)
	item2 := g.NewObj(invTestEquip)
	item3 := g.NewObj(invTestEquip)

	inv.Add(item1)
	inv.Add(item2)
	inv.Add(item3)

	taken := inv.Take(1)

	if taken != item2 {
		t.Errorf(`inv.Take(1) gave %v, want %v`, taken, item2)
	}

	if size := inv.Len(); size != 2 {
		t.Errorf(`inv.Len() was %v after taking, want 2`, size)
	}

	taken = inv.Take(1)

	if taken != item3 {
		t.Errorf(`inv.Take(1) gave %v, want %v`, taken, item3)
	}

	if size := inv.Len(); size != 1 {
		t.Errorf(`inv.Len() was %v after taking, want 2`, size)
	}
}

func TestAddOverCapacity(t *testing.T) {
	g := newTestGame()
	inv := NewInventoryWithCap(2)

	item1 := g.NewObj(invTestEquip)
	item2 := g.NewObj(invTestEquip)
	item3 := g.NewObj(invTestEquip)

	inv.Add(item1)
	inv.Add(item2)

	if inv.Add(item3) {
		t.Errorf(`inv.Add(item) over capacity was true, want false`)
	}

	if inv.Len() > 2 {
		t.Errorf(`inv.Add(item) over capacity 2 still added past 2`)
	}
}

func TestHasEquips(t *testing.T) {
	g := newTestGame()
	inv := NewInventoryWithCap(2)

	if inv.HasEquipment() {
		t.Errorf(`inv.HasEquips() was true, want false`)
	}

	inv.Add(g.NewObj(invTestConsumable))

	if inv.HasEquipment() {
		t.Errorf(`inv.HasEquips() was true, want false`)
	}

	inv.Add(g.NewObj(invTestEquip))

	if !inv.HasEquipment() {
		t.Errorf(`inv.HasEquips() was false, want true`)
	}
}

func TestHasUsables(t *testing.T) {
	g := newTestGame()
	inv := NewInventoryWithCap(2)

	if inv.HasUsables() {
		t.Errorf(`inv.HasUsables() was true, want false`)
	}

	inv.Add(g.NewObj(invTestEquip))

	if inv.HasUsables() {
		t.Errorf(`inv.HasUsables() was true, want false`)
	}

	inv.Add(g.NewObj(invTestConsumable))

	if !inv.HasUsables() {
		t.Errorf(`inv.HasUsables() was false, want true`)
	}
}
