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
	slot := item.Equipment.Slot
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
	for _, slot := range b.Slots {
		if slot != nil {
			return false
		}
	}
	return true
}

// Get the total bonus/malus to melee from equipment worn on this body.
func (b *Body) Melee() int {
	melee := 0
	for _, equip := range b.all() {
		melee += equip.Equipment.Melee
	}
	return melee
}

// Get the total bonus/malus to evasion from equipment worn on this body.
func (b *Body) Evasion() int {
	evasion := 0
	for _, equip := range b.all() {
		evasion += equip.Equipment.Evasion
	}
	return evasion
}

func (b *Body) ProtDice() []Dice {
	dice := make([]Dice, 0, numSlots)
	for _, equip := range b.all() {
		dice = append(dice, equip.Equipment.Protroll)
	}
	return dice
}

func (b *Body) Weapon() *Obj {
	return b.Slots[SlotHand]
}

// Accumulate all the effects on all of our armor. The result is a histogram of
// how often each effect appeared on this body's collection of armor.
func (b *Body) ArmorEffects() map[Effect]int {
	effects := make(map[Effect]int)

	for slot, equip := range b.all() {
		if slot == SlotHand {
			continue
		}
		for _, e := range equip.Equipment.Effects {
			effects[e] += 1
		}
	}
	return effects
}

// Return a collection of all the equipped stuff on this body, without all the
// nil slots.
func (b *Body) all() []*Obj {
	equips := make([]*Obj, 0, numSlots)
	for _, equip := range b.Slots {
		if equip != nil {
			equips = append(equips, equip)
		}
	}
	return equips
}
