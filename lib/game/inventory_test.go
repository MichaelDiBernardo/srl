package game

import (
	"testing"
)

var invTestQueue = newEventQueue()

var invTestItem = &Spec{
	Family:  FamItem,
	Genus:   GenEquip,
	Species: "testspec",
	Name:    "Item",
	Traits:  &Traits{},
}

func TestTake(t *testing.T) {
	inv := NewInventory()

	item1 := newObj(invTestItem, invTestQueue)
	item2 := newObj(invTestItem, invTestQueue)
	item3 := newObj(invTestItem, invTestQueue)

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
