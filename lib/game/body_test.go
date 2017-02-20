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

func TestArmorEffects(t *testing.T) {
	// Sorry I was too lazy to register "testing effects"
	const (
		r1 = BrandFire
		r2 = BrandElec
		r3 = BrandIce
	)
	var (
		a1 = &Spec{
			Family:  FamItem,
			Genus:   GenEquipment,
			Species: "testspec",
			Name:    "head",
			Traits: &Traits{
				Equipment: NewEquipment(Equipment{
					Slot:    SlotHead,
					Effects: NewEffects(map[Effect]int{r1: 1}),
				}),
			},
		}
		a2 = &Spec{
			Family:  FamItem,
			Genus:   GenEquipment,
			Species: "testspec",
			Name:    "arms",
			Traits: &Traits{
				Equipment: NewEquipment(Equipment{
					Slot:    SlotArms,
					Effects: NewEffects(map[Effect]int{r2: 1}),
				}),
			},
		}
		a3 = &Spec{
			Family:  FamItem,
			Genus:   GenEquipment,
			Species: "testspec",
			Name:    "legs",
			Traits: &Traits{
				Equipment: NewEquipment(Equipment{
					Slot:    SlotLegs,
					Effects: NewEffects(map[Effect]int{r2: 1, r3: 1}),
				}),
			},
		}
	)

	b := NewBody()
	b.Wear(newObj(a1))
	b.Wear(newObj(a2))
	b.Wear(newObj(a3))

	effects := b.ArmorEffects()

	if n := effects[r1].Count; n != 1 {
		t.Errorf(`ArmorEffects() returned effects[r1].Count = %d; want 1 %v`, n, effects)
	}
	if n := effects[r2].Count; n != 2 {
		t.Errorf(`ArmorEffects() returned effects[r2].Count = %d; want 2`, n)
	}
	if n := effects[r3].Count; n != 1 {
		t.Errorf(`ArmorEffects() returned effects[r3].Count = %d; want 1`, n)
	}
}
