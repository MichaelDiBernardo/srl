package game

import (
	"testing"
)

var btWeapSpec = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
	Species: "testspec",
	Name:    "Hand",
	Traits: &Traits{
		Equipment: NewEquipment(Equipment{Slot: SlotHand}),
	},
}

func TestWearIntoEmptySlot(t *testing.T) {
	b := NewBody()
	weap := newObj(btWeapSpec)
	b.Wear(weap)

	if w := b.Slots[SlotHand]; w != weap {
		t.Errorf(`Found %p in hand, want %p`, w, weap)
	}
}

func TestWearIntoOccupiedSlot(t *testing.T) {
	b := NewBody()
	weap1, weap2 := newObj(btWeapSpec), newObj(btWeapSpec)
	b.Wear(weap1)

	if old := b.Wear(weap2); old != weap1 {
		t.Errorf(`Returned %p from hand, want %p`, old, weap1)
	}

	if w := b.Slots[SlotHand]; w != weap2 {
		t.Errorf(`Found %p in hand, want %p`, w, weap2)
	}
}

func TestRemoveFromEmptySlotReturnsNil(t *testing.T) {
	b := NewBody()

	if w := b.Remove(SlotHand); w != nil {
		t.Errorf(`Removed %p, want nil`, w)
	}
}

func TestRemoveFromFullSlotReturnsItem(t *testing.T) {
	b := NewBody()
	weap := newObj(btWeapSpec)
	b.Wear(weap)

	if w := b.Remove(SlotHand); w != weap {
		t.Errorf(`Removed %p, want %p`, w, weap)
	}

	if w := b.Slots[SlotHand]; w != nil {
		t.Errorf(`Left %p behind, want nil`, w)
	}
}

func TestIsNaked(t *testing.T) {
	b := NewBody()
	if b.Naked() == false {
		t.Error(`b.Naked() was false, want true`)
	}
}

func TestIsNotNaked(t *testing.T) {
	b := NewBody()
	b.Wear(newObj(btWeapSpec))

	if b.Naked() == true {
		t.Error(`b.Naked() was true, want false`)
	}
}
