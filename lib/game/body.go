package game

import ()

// Slots where items are worn/wielded on an actor's body.
type Slot int

const (
	SlotHand = iota
	SlotOffHand
	numSlots
)

type Body struct {
	Slots [numSlots]*Obj
}

func (b *Body) Wear(item *Obj) *Obj {
	slot := item.Equip.Slot
	equipped := b.Slots[slot]
	b.Slots[slot] = item
	return equipped
}

func (b *Body) Remove(slot Slot) *Obj {
	equipped := b.Slots[slot]
	b.Slots[slot] = nil
	return equipped
}
