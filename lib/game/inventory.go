package game

import (
	"container/list"
	"fmt"
)

// A container that holds any type of item.
type Inventory struct {
	Items *list.List
}

// Create a new empty inventory.
func NewInventory() *Inventory {
	return &Inventory{Items: list.New()}
}

// The number of items in this inventory.
func (inv *Inventory) Len() int {
	return inv.Items.Len()
}

// Is this inventory empty?
func (inv *Inventory) Empty() bool {
	return inv.Items.Len() == 0
}

// Tries to add item to this inventory. Returns false if the item doesn't fit.
func (inv *Inventory) Add(item *Obj) bool {
	if fam := item.Spec.Family; fam != FamItem {
		panic(fmt.Sprintf("Tried to add obj of family %v to inventory.", fam))
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

// Returns the item at index 'index' and removes it from the inventory.
func (inv *Inventory) Take(index int) *Obj {
	if inv.Empty() {
		panic(fmt.Sprintf(`Tried to Take(%v) from empty inventory.`, index))
	}
	if size := inv.Len(); size <= index {
		panic(fmt.Sprintf(`Tried to Take(%v) from inventory of size %v.`, index, size))
	}

	itemElem := inv.Items.Back()
	for i := 0; i != index; i++ {
		itemElem = itemElem.Prev()
	}

	inv.Items.Remove(itemElem)
	return itemElem.Value.(*Obj)
}
