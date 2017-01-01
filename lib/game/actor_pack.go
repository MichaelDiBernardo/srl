package game

import (
	"fmt"
)

// A thing that that can hold items in inventory. (A "pack".)
type Packer interface {
	Objgetter
	// Tries to pickup something at current square. If there are many things,
	// will invoke stack menu.
	TryPickup()
	// Pickup the item on the floor stack at given index.
	Pickup(index int)
	// Tries to drop something at current square.
	TryDrop()
	// Drop the item at index in inventory to the floor stack.
	Drop(index int)
	// Get this Packer's inventory.
	Inventory() *Inventory
}

// An attacker that works for all actors.
type ActorPacker struct {
	Trait
	inventory *Inventory
}

func NewActorPacker(obj *Obj) Packer {
	return &ActorPacker{
		Trait:     Trait{obj: obj},
		inventory: NewInventory(),
	}
}

func (a *ActorPacker) Inventory() *Inventory {
	return a.inventory
}

func (a *ActorPacker) TryPickup() {
	ground := a.obj.Tile.Items
	if ground.Empty() {
		a.obj.Game.Events.Message("Nothing there.")
	} else if ground.Len() == 1 {
		a.moveFromGround(0)
	} else {
		a.obj.Game.SwitchMode(ModePickup)
	}
}

func (a *ActorPacker) Pickup(index int) {
	a.obj.Game.SwitchMode(ModeHud)
	a.moveFromGround(index)
}

func (a *ActorPacker) TryDrop() {
	if a.inventory.Empty() {
		a.obj.Game.Events.Message("Nothing to drop.")
		return
	}

	ground := a.obj.Tile.Items
	if ground.Full() {
		a.obj.Game.Events.Message("Can't drop here.")
	} else {
		a.obj.Game.SwitchMode(ModeDrop)
	}
}

func (a *ActorPacker) Drop(index int) {
	a.obj.Game.SwitchMode(ModeHud)
	item := a.inventory.Take(index)

	// Bounds-check the index the player requested.
	if item == nil {
		return
	}
	a.obj.Tile.Items.Add(item)
	a.obj.Game.Events.Message(fmt.Sprintf("%v dropped %v.", a.obj.Spec.Name, item.Spec.Name))
}

func (a *ActorPacker) moveFromGround(index int) {
	// Bounds-check the index the player requested.
	item := a.obj.Tile.Items.At(index)
	if item == nil {
		return
	}

	if a.inventory.Full() {
		a.obj.Game.Events.Message(fmt.Sprintf("%v has no room for %v.", a.obj.Spec.Name, item.Spec.Name))
	} else {
		item := a.obj.Tile.Items.Take(index)
		a.inventory.Add(item)
		a.obj.Game.Events.Message(fmt.Sprintf("%v got %v.", a.obj.Spec.Name, item.Spec.Name))
	}
}
