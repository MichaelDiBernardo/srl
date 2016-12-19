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

// Get the total bonus/malus to melee from equipment worn on this body.
func (b *Body) Melee() int {
	melee := 0
	equips := b.all()
	for i := 0; i < len(equips); i++ {
		melee += equips[i].Equip.Melee
	}
	return melee
}

// Get the total bonus/malus to evasion from equipment worn on this body.
func (b *Body) Evasion() int {
	evasion := 0
	equips := b.all()
	for i := 0; i < len(equips); i++ {
		evasion += equips[i].Equip.Evasion
	}
	return evasion
}

func (b *Body) ProtDice() []Dice {
	dice := make([]Dice, 0, numSlots)
	equips := b.all()
	for i := 0; i < len(equips); i++ {
		dice = append(dice, equips[i].Equip.Protroll)
	}
	return dice
}

func (b *Body) Weapon() *Obj {
	return b.Slots[SlotHand]
}

// Return a collection of all the equipped stuff on this body, without all the
// nil slots.
func (b *Body) all() []*Obj {
	equips := make([]*Obj, 0, numSlots)
	for i := 0; i < numSlots; i++ {
		e := b.Slots[i]
		if e != nil {
			equips = append(equips, e)
		}
	}
	return equips
}
