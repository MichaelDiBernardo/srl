package game

import (
	"fmt"
)

// A thing that that can hold items in inventory. (A "pack".)
type Packer interface {
	Objgetter
	// Tries to pickup something at current square. If there are many things,
	// will invoke stack menu. Returns true if a turn should pass because the
	// player picked up a single item below them without switching modes.
	TryPickup() bool
	// Pickup the item on the floor stack at given index. Return true if a turn
	// should pass.
	Pickup(index int) bool
	// Tries to drop something at current square.
	TryDrop()
	// Drop the item at index in inventory to the floor stack. Return true if
	// an action was taken that requires a turn to pass.
	Drop(index int) bool
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

func (a *ActorPacker) TryPickup() bool {
	evolve := false
	ground := a.obj.Tile.Items
	if ground.Empty() {
		a.obj.Game.Events.Message("Nothing there.")
	} else if ground.Len() == 1 {
		evolve = a.moveFromGround(0)
	} else {
		a.obj.Game.SwitchMode(ModePickup)
	}
	return evolve
}

func (a *ActorPacker) Pickup(index int) bool {
	a.obj.Game.SwitchMode(ModeHud)
	return a.moveFromGround(index)
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

// Returns false if no actual action was taken.
func (a *ActorPacker) Drop(index int) bool {
	a.obj.Game.SwitchMode(ModeHud)
	item := a.inventory.Take(index)

	// Bounds-check the index the player requested.
	if item == nil {
		return false
	}

	a.obj.Tile.Items.Add(item)
	a.obj.Game.Events.Message(fmt.Sprintf("%v dropped %v.", a.obj.Spec.Name, item.Spec.Name))

	return true
}

func (a *ActorPacker) moveFromGround(index int) bool {
	// Bounds-check the index the player requested.
	item := a.obj.Tile.Items.At(index)
	if item == nil {
		return false
	}

	if a.inventory.Full() {
		a.obj.Game.Events.Message(fmt.Sprintf("%v has no room for %v.", a.obj.Spec.Name, item.Spec.Name))
		return false
	}

	item = a.obj.Tile.Items.Take(index)
	a.inventory.Add(item)
	a.obj.Game.Events.Message(fmt.Sprintf("%v got %v.", a.obj.Spec.Name, item.Spec.Name))
	return true
}
