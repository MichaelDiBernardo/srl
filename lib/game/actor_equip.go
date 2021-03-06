package game

import (
	"fmt"
)

type Equipper interface {
	Objgetter
	// Bring up the equipper screen if anything in inventory can be equipped.
	TryEquip()
	// Bring up the remover screen if anything on body can be removed.
	TryRemove()
	// Equip the item at index 'index' in inventory. Return true if a turn
	// should pass.
	Equip(index int) bool
	// Remove the item equipped in the given slot. Return true if a turn should
	// pass.
	Remove(slot Slot) bool
	// Get the underlying entity's Body.
	Body() *Body
}

type ActorEquipper struct {
	Trait
	body *Body
}

func NewActorEquipper(obj *Obj) Equipper {
	return &ActorEquipper{
		Trait: Trait{obj: obj},
		body:  NewBody(),
	}
}

func (a *ActorEquipper) TryEquip() {
	if !a.obj.Packer.Inventory().HasEquipment() {
		a.obj.Game.Events.Message("Nothing to wield/wear.")
	} else {
		a.obj.Game.SwitchMode(ModeEquip)
	}
}

func (a *ActorEquipper) TryRemove() {
	if a.body.Naked() {
		a.obj.Game.Events.Message("Not wearing anything.")
	} else if a.obj.Packer.Inventory().Full() && a.obj.Tile.Items.Full() {
		a.obj.Game.Events.Message("Can't remove; pack and ground are full.")
	} else {
		a.obj.Game.SwitchMode(ModeRemove)
	}
}

func (a *ActorEquipper) Equip(index int) bool {
	a.obj.Game.SwitchMode(ModeHud)
	inv := a.obj.Packer.Inventory()

	equip := inv.At(index)

	// Bounds-check the index the player requested.
	if equip == nil {
		return false
	}

	if equip.Spec.Genus != GenEquipment {
		a.obj.Game.Events.Message(fmt.Sprintf("Cannot equip %v.", equip.Spec.Name))
		return false
	}
	equip = inv.Take(index)

	if swapped := a.body.Wear(equip); swapped != nil {
		a.obj.Packer.Inventory().Add(swapped)
	}
	return true
}

func (a *ActorEquipper) Remove(slot Slot) bool {
	a.obj.Game.SwitchMode(ModeHud)

	removed := a.body.Remove(slot)
	if removed == nil {
		return false
	}

	if added := a.obj.Packer.Inventory().Add(removed); added {
		return true
	}

	// No room for unequipped item in inventory; drop it.
	a.obj.Tile.Items.Add(removed)
	a.obj.Game.Events.Message(fmt.Sprintf("No room in pack! Dropped %v.", removed.Spec.Name))
	return true
}

func (a *ActorEquipper) Body() *Body {
	return a.body
}
