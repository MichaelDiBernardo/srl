package game

import (
	"fmt"
)

// A thing that can use single-use things.
type User interface {
	Objgetter
	// Bring up the 'use' screen if anything in inventory can be used.
	TryUse()
	// Use the item at index 'index' in inventory. Return true if a turn should
	// pass.
	Use(index int) bool
}

type ActorUser struct {
	Trait
}

func NewActorUser(obj *Obj) User {
	return &ActorUser{Trait: Trait{obj: obj}}
}

func (a *ActorUser) TryUse() {
	if !a.obj.Packer.Inventory().HasUsables() {
		a.obj.Game.Events.Message("Nothing to use.")
	} else {
		a.obj.Game.SwitchMode(ModeUse)
	}
}

func (a *ActorUser) Use(index int) bool {
	a.obj.Game.SwitchMode(ModeHud)
	inv := a.obj.Packer.Inventory()

	item := inv.At(index)

	// Bounds-check the index the player requested.
	if item == nil {
		return false
	}

	if item.Spec.Genus != GenConsumable {
		a.obj.Game.Events.Message(fmt.Sprintf("Cannot use %v.", item.Spec.Name))
		return false
	}

	item = inv.Take(index)
	item.Consumable.Consume(a)
	return true
}
