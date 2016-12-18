package game

import ()

// Slots where items are worn/wielded on an actor's body.
type Slot int

const (
	SlotHand = iota
	SlotHead
	SlotBody
	SlotArms
	SlotLegs
	SlotRelic
	numSlots
)

type Body struct {
	Slots [numSlots]*Obj
}

func NewBody() *Body {
	return &Body{}
}

func (b *Body) Wear(item *Obj) *Obj {
	slot := item.Equip.Slot
	equipped := b.Slots[slot]
	b.Slots[slot] = item
	return equipped
}

func (b *Body) Remove(slot Slot) *Obj {
	if slot >= numSlots {
		return nil
	}
	equipped := b.Slots[slot]
	b.Slots[slot] = nil
	return equipped
}

func (b *Body) Naked() bool {
	for i := 0; i < numSlots; i++ {
		if b.Slots[i] != nil {
			return false
		}
	}
	return true
}
