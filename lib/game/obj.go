package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
)

// e.g. "Actor", "Item"
type Family string

// e.g. "Player", "Monster", "Potion"
type Genus string

// e.g. "Orc", "HealPotion"
type Species string

// Permissible families of objects.
const (
	FamActor = "actor"
	FamItem  = "item"
)

// Trait basics; the ability to backreference the attached object.
type Objgetter interface {
	Obj() *Obj
}

type Trait struct {
	obj *Obj
}

func (t *Trait) Obj() *Obj {
	return t.obj
}

// Stores everything you need to know about randomly generating this object
// ingame.
// 'Depths' is a list of the native depths of this object.
// 'GroupSize' means "pack size" for monsters and "stack size" for consumables.
type Gen struct {
	Depths    []int
	GroupSize int
}

// Should this entry be "findable" in the given range of depths?
func (g Gen) Findable(low, high int) bool {
	log.Printf("Gen.Findable: checking [%d,%d]", low, high)
	for _, d := range g.Depths {
		if low <= d && d <= high {
			log.Printf("%d: good!", d)
			return true
		}
	}
	log.Printf("Gen.Findable: no good depths for [%d,%d]", low, high)
	return false
}

// A specification for a type of game object.
type Spec struct {
	Family  Family
	Genus   Genus
	Species Species
	Name    string
	Traits  *Traits
	Gen     Gen
}

var nextobjid = 1

// Specifically, an in-game object that can be placed on a map and can Do
// Something. Its traits determine what it can do.
type Obj struct {
	id   int
	Spec *Spec

	Tile  *Tile
	Level *Level
	Game  *Game

	// Actor traits.
	Mover    Mover
	AI       AI
	Sheet    Sheet
	Fighter  Fighter
	Packer   Packer
	Equipper Equipper
	User     User

	// Item traits. Since these don't ever conceivably need alternate
	// implementations, they are not interface types.
	Equipment  *Equipment
	Consumable *Consumable
}

func (o *Obj) String() string {
	if o.Spec != nil {
		return "Obj: " + o.Spec.Name
	}
	return "Obj: ???"
}

// A specification object for newObj. Each key maps to a factory function for
// the specific implementation of the desired trait. If an object is not
// supposed to have a specific trait, leave it unspecified.
type Traits struct {
	Mover    func(*Obj) Mover
	AI       func(*Obj) AI
	Sheet    func(*Obj) Sheet
	Fighter  func(*Obj) Fighter
	Packer   func(*Obj) Packer
	Equipper func(*Obj) Equipper
	User     func(*Obj) User

	Equipment  func(*Obj) *Equipment
	Consumable func(*Obj) *Consumable
}

// Create a new game object from the given spec. This shouldn't be used
// directly; you should instead use a *Game as a factory for any game objects
// that need creating. This will not initialize the fields on the obj that have
// nothing to do with specs or traits (e.g. game, eventqueue, tile etc.)
func newObj(spec *Spec) *Obj {
	// Create.
	newobj := &Obj{Spec: spec, id: nextobjid}
	nextobjid++

	// Assign traits.
	traits := spec.Traits
	if traits.Mover != nil {
		newobj.Mover = traits.Mover(newobj)
	}
	if traits.AI != nil {
		newobj.AI = traits.AI(newobj)
	}
	if traits.Sheet != nil {
		newobj.Sheet = traits.Sheet(newobj)
	}
	if traits.Fighter != nil {
		newobj.Fighter = traits.Fighter(newobj)
	}
	if traits.Packer != nil {
		newobj.Packer = traits.Packer(newobj)
	}
	if traits.Equipper != nil {
		newobj.Equipper = traits.Equipper(newobj)
	}
	if traits.User != nil {
		newobj.User = traits.User(newobj)
	}

	if traits.Equipment != nil {
		newobj.Equipment = traits.Equipment(newobj)
	}
	if traits.Consumable != nil {
		newobj.Consumable = traits.Consumable(newobj)
	}
	return newobj
}

// What point on the map is this object on?
func (o *Obj) Pos() math.Point {
	return o.Tile.Pos
}

// Does this object represent the player?
func (o *Obj) IsPlayer() bool {
	return o.Spec.Family == FamActor && o.Spec.Genus == GenPlayer
}

func (o *Obj) HasAI() bool {
	return o.AI != nil
}
