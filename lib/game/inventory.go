package game

import (
	"container/list"
	"fmt"
)

const inventoryCapacity = 20

// A container that holds any type of item.
type Inventory struct {
	Items    *list.List
	capacity int
}

// Create a new empty inventory with capacity `inventoryCapacity`.
func NewInventory() *Inventory {
	return &Inventory{Items: list.New(), capacity: inventoryCapacity}
}

func NewInventoryWithCap(capacity int) *Inventory {
	return &Inventory{Items: list.New(), capacity: capacity}
}

// The number of items in this inventory.
func (inv *Inventory) Len() int {
	return inv.Items.Len()
}

// Is this inventory empty?
func (inv *Inventory) Empty() bool {
	return inv.Items.Len() == 0
}

// Is this inventory at capacity?
func (inv *Inventory) Full() bool {
	return inv.Items.Len() >= inv.capacity
}

// Tries to add item to this inventory. Returns false if the item doesn't fit.
func (inv *Inventory) Add(item *Obj) bool {
	if fam := item.Spec.Family; fam != FamItem {
		panic(fmt.Sprintf("Tried to add obj of family %v to inventory.", fam))
	}
	if inv.Full() {
		return false
	}
	inv.Items.PushFront(item)
	return true
}

// Gets the item at the "top" of this inventory. This item can be used to
// visually represent the entire inventory if it is a floor stack, for example.
// Returns nil if inventory is empty.
func (inv *Inventory) Top() *Obj {
	if inv.Empty() {
		return nil
	}
	return inv.Items.Front().Value.(*Obj)
}

// Peeks at the item at the given index without removing it.
// Returns nil if there is no item at that index.
func (inv *Inventory) At(index int) *Obj {
	if size := inv.Len(); size <= index {
		return nil
	}
	return inv.itemElemAt(index).Value.(*Obj)
}

// Returns the item at index 'index' and removes it from the inventory.
// Returns nil if there was no item at the given index.
func (inv *Inventory) Take(index int) *Obj {
	if size := inv.Len(); size <= index {
		return nil
	}

	itemElem := inv.itemElemAt(index)
	inv.Items.Remove(itemElem)
	return itemElem.Value.(*Obj)
}

// Does this inventory have anything to equip in it?
func (inv *Inventory) HasEquips() bool {
	for e := inv.Items.Front(); e != nil; e = e.Next() {
		item := e.Value.(*Obj)
		if item.Spec.Genus == GenEquip {
			return true
		}
	}
	return false
}

func (inv *Inventory) itemElemAt(index int) *list.Element {
	itemElem := inv.Items.Back()
	for i := 0; i != index; i++ {
		itemElem = itemElem.Prev()
	}
	return itemElem
}
