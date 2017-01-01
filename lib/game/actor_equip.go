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
	// Equip the item at index 'index' in inventory.
	Equip(index int)
	// Remove the item equipped in the given slot.
	Remove(slot Slot)
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

func (a *ActorEquipper) Equip(index int) {
	a.obj.Game.SwitchMode(ModeHud)

	equip := a.obj.Packer.Inventory().Take(index)

	// Bounds-check the index the player requested.
	if equip == nil {
		return
	}

	if equip.Spec.Genus != GenEquip {
		a.obj.Game.Events.Message(fmt.Sprintf("Cannot equip %v.", equip.Spec.Name))
		return
	}

	if swapped := a.body.Wear(equip); swapped != nil {
		a.obj.Packer.Inventory().Add(swapped)
	}
}

func (a *ActorEquipper) Remove(slot Slot) {
	a.obj.Game.SwitchMode(ModeHud)

	removed := a.body.Remove(slot)
	if removed == nil {
		return
	}

	if added := a.obj.Packer.Inventory().Add(removed); added {
		return
	}

	// No room for unequipped item in inventory; drop it.
	a.obj.Tile.Items.Add(removed)
	a.obj.Game.Events.Message(fmt.Sprintf("No room in pack! Dropped %v.", removed.Spec.Name))
}

func (a *ActorEquipper) Body() *Body {
	return a.body
}
