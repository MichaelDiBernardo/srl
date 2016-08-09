package game

const GenEquip = "equip"

type Slot int

const (
	SlotHand = iota
	SlotOffHand
)

type Equip interface {
	Objgetter
	Slot() Slot
	Desc() string
	// Equip getters. This won't have a lot of behaviour because all the combat logic should go in Fighter.
}
