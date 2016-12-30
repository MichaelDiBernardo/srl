package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

var (
	atActorSpec = &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Name:    "Hi",
		Traits: &Traits{
			Mover:    NewActorMover,
			Packer:   NewActorPacker,
			Equipper: NewActorEquipper,
		},
	}

	atItemSpec = &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: "testspec",
		Name:    "Item",
		Traits: &Traits{
			Equip: NewEquip(Equip{Slot: SlotHand}),
		},
	}
)

func TestOkMove(t *testing.T) {
	g := NewGame()
	l := NewLevel(4, 4, g, IdentLevel)
	obj := g.NewObj(atActorSpec)
	startpos := math.Pt(1, 1)

	l.Place(obj, startpos)

	ok := obj.Mover.Move(math.Pt(1, 0))

	if !ok {
		t.Error(`Move( (1, 0)) = false, want true`)
	}

	newpos := obj.Pos()
	want := math.Pt(2, 1)
	if newpos != want {
		t.Errorf(`Move((1, 0)) = %v, want %v`, newpos, want)
	}

	if l.At(startpos).Actor != nil {
		t.Error(`Move((1, 0)) did not set start tile actor to nil`)
	}
	if l.At(newpos).Actor != obj {
		t.Error(`Move((1, 0)) did not set dest tile actor to obj`)
	}
}

func TestActorCollision(t *testing.T) {
	g := NewGame()
	l := NewLevel(4, 4, g, IdentLevel)
	a1, a2 := g.NewObj(atActorSpec), g.NewObj(atActorSpec)
	l.Place(a1, math.Pt(1, 1))
	l.Place(a2, math.Pt(2, 1))

	ok := a1.Mover.Move(math.Pt(1, 0))

	if ok {
		t.Error(`a1.Move( (1, 0)) = true, want false`)
	}
}

func TestPlayerMaxHPCalc(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, vit: 1}

	if maxhp, want := obj.Sheet.MaxHP(), 20; maxhp != want {
		t.Errorf(`MaxHP() was %d, want %d`, maxhp, want)
	}
}

func TestPlayerMaxMPCalc(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, mnd: 2}

	if maxmp, want := obj.Sheet.MaxMP(), 30; maxmp != want {
		t.Errorf(`MaxMP() was %d, want %d`, maxmp, want)
	}
}

func TestHurtingPlayerToDeathEndsGame(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 0, agi: 0}
	g.Player = obj

	obj.Sheet.Hurt(9999999)
	if m := g.mode; m != ModeGameOver {
		t.Errorf(`Killing player changed mode to %v; want %v`, m, ModeGameOver)
	}
}

var knifeSpec = &Spec{
	Family:  FamItem,
	Genus:   GenEquip,
	Species: SpecSword,
	Name:    "KNIFE",
	Traits: &Traits{
		Equip: NewEquip(Equip{
			Damroll: NewDice(1, 7),
			Melee:   1,
			Evasion: 1,
			Weight:  2,
			Slot:    SlotHand,
		}),
	},
}

func testAtkEq(t *testing.T, atk Attack, want Attack) {
	if m, w := atk.Melee, want.Melee; m != w {
		t.Errorf(`atk.Melee was %d, want %d`, m, w)
	}
	if d, w := atk.Damroll.Dice, want.Damroll.Dice; d != w {
		t.Errorf(`atk.Damroll.Dice was %d, want %d`, d, w)
	}
	if s, w := atk.Damroll.Sides, want.Damroll.Sides; s != w {
		t.Errorf(`atk.Damroll.Sides was %d, want %d`, s, w)
	}
}

func TestPlayerAttackNoBonuses(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 0, agi: 0}

	weap := g.NewObj(knifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 7)})
}

func TestPlayerAttackStrBonusBelowCap(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 1, agi: 0}

	weap := g.NewObj(knifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 8)})
}

func TestPlayerAttackStrBonusAboveCap(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 3, agi: 0}

	weap := g.NewObj(knifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 9)})
}

func TestPlayerAttackMeleeBonus(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 0, agi: 3}

	weap := g.NewObj(knifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 4, Damroll: NewDice(1, 7)})
}

func TestPlayerAttackFistNoStrSides(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 0, agi: 0}

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 0, Damroll: NewDice(1, 1)})
}

func TestPlayerAttackFistStrSides(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 10, agi: 0}

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 0, Damroll: NewDice(1, 11)})
}

func testDefEq(t *testing.T, def Defense, want Defense) {
	if e, w := def.Evasion, want.Evasion; e != w {
		t.Errorf(`def.Evasion was %d, want %d`, e, w)
	}
	if l, w := len(def.ProtDice), len(want.ProtDice); l != w {
		t.Errorf(`len(def.ProtDice) was %d, want %d`, l, w)
	}
	for i, d := range def.ProtDice {
		w := want.ProtDice[i]
		if d != w {
			t.Errorf(`def.ProtDice[%d] was '%s', want '%s'`, i, d, w)
		}
	}
}

func TestPlayerDefenseNoArmorOrEvasion(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, agi: 0}

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{Evasion: 0})
}

func TestPlayerDefenseNoArmorWithEvasion(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, agi: 2}

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{Evasion: 2})
}

func TestPlayerDefenseWithArmor(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, agi: 0}

	armspec1 := &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: SpecLeatherArmor,
		Name:    "LEATHER",
		Traits: &Traits{
			Equip: NewEquip(Equip{
				Protroll: NewDice(1, 4),
				Melee:    0,
				Evasion:  -1,
				Weight:   4,
				Slot:     SlotBody,
			}),
		},
	}
	armspec2 := &Spec{
		Family:  FamItem,
		Genus:   GenEquip,
		Species: SpecLeatherArmor,
		Name:    "MASK",
		Traits: &Traits{
			Equip: NewEquip(Equip{
				Protroll: NewDice(1, 3),
				Melee:    0,
				Evasion:  -2,
				Weight:   2,
				Slot:     SlotHead,
			}),
		},
	}

	obj.Equipper.Body().Wear(g.NewObj(armspec1))
	obj.Equipper.Body().Wear(g.NewObj(armspec2))

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{
		Evasion: -3,
		ProtDice: []Dice{
			NewDice(1, 3),
			NewDice(1, 4),
		},
	})
}

type fakefighter struct {
	Trait
	Called bool
}

func (f *fakefighter) Hit(other Fighter) {
	f.Called = true
}

func TestPlayerMonsterCollisionsHit(t *testing.T) {
	g := NewGame()
	player := g.NewObj(PlayerSpec)
	pf := &fakefighter{Trait: Trait{obj: player}}
	player.Fighter = pf

	monster := g.NewObj(atActorSpec)
	mf := &fakefighter{Trait: Trait{obj: player}}
	monster.Fighter = mf

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(player, math.Pt(0, 0))
	l.Place(monster, math.Pt(1, 1))

	player.Mover.Move(math.Pt(1, 1))

	if !pf.Called {
		t.Error("Moving player into other did not try to hit.")
	}

	monster.Mover.Move(math.Pt(-1, -1))

	if !mf.Called {
		t.Error("Moving other into player did not try to hit.")
	}
}

func TestMonsterMonsterCollisionsHit(t *testing.T) {
	g := NewGame()
	mon1 := g.NewObj(atActorSpec)
	mf1 := &fakefighter{Trait: Trait{obj: mon1}}
	mon1.Fighter = mf1

	mon2 := g.NewObj(atActorSpec)
	mf2 := &fakefighter{Trait: Trait{obj: mon2}}
	mon2.Fighter = mf2

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(mon1, math.Pt(0, 0))
	l.Place(mon2, math.Pt(1, 1))

	mon1.Mover.Move(math.Pt(1, 1))

	if mf1.Called {
		t.Error("Moving monster into monster tried to hit.")
	}
}

func TestTryPickupNoItemsOnGround(t *testing.T) {
	g := NewGame()
	taker := g.NewObj(atActorSpec)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(taker, math.Pt(0, 0))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size > 0 {
		t.Errorf(`TryPickup() on empty square gave inven size %d; want 0`, size)
	}
}

func TestTryPickupOneItemOnGround(t *testing.T) {
	g := NewGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(taker, math.Pt(0, 0))
	l.Place(item, math.Pt(0, 0))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size != 1 {
		t.Errorf(`TryPickup() on 1-item square gave inven size %d; want 1`, size)
	}
	if size := l.At(math.Pt(0, 0)).Items.Len(); size != 0 {
		t.Errorf(`TryPickup() on 1-item square left %d items on ground; want 0`, size)
	}
}

func TestTryPickupFromStack(t *testing.T) {
	g := NewGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)
	item2 := g.NewObj(atItemSpec)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(taker, math.Pt(0, 0))
	l.Place(item, math.Pt(0, 0))
	l.Place(item2, math.Pt(0, 0))

	taker.Packer.TryPickup()
	if size := taker.Packer.Inventory().Len(); size != 0 {
		t.Errorf(`TryPickup() on stack took something instead of opening menu; took %d things`, size)
	}
	if size := l.At(math.Pt(0, 0)).Items.Len(); size != 2 {
		t.Errorf(`TryPickup() took from ground instead of opening menu; left %d things`, size)
	}
	if size := g.Events.Len(); size != 1 {
		t.Errorf(`TryPickup() pushed wrong # of events to queue; found %d, want 1`, size)
	}

	e, ok := g.Events.Next().(ModeEvent)
	if !ok {
		t.Error(`TryPickup pushed wrong event type to queue.`)
	}
	if e.Mode != ModePickup {
		t.Errorf(`TryPickup switched to mode %v, want %v`, e.Mode, ModePickup)
	}

}

func TestPickupOutOfBounds(t *testing.T) {
	g := NewGame()
	taker := g.NewObj(atActorSpec)
	item := g.NewObj(atItemSpec)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(taker, math.Pt(0, 0))
	l.Place(item, math.Pt(0, 0))

	taker.Packer.TryPickup()
	taker.Packer.Pickup(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Out-of-bounds Pickup switched to mode %v; want %v`, mode, ModeHud)
	}

}

func TestTryEquipWithNoEquipsInInventory(t *testing.T) {
	g := NewGame()
	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()
	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryEquip w no equips switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryEquipWithEquipsInInventory(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	equipper.Packer.Inventory().Add(equip)

	equipper.Equipper.TryEquip()

	if mode := g.mode; mode != ModeEquip {
		t.Errorf(`TryEquip switched to mode %v, want %v`, mode, ModeEquip)
	}
}

func TestEquipIntoEmptySlot(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	inv := equipper.Packer.Inventory()
	inv.Add(equip)

	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Was mode %v after equip; want %v`, mode, ModeHud)
	}

	if !inv.Empty() {
		t.Errorf(`Item did not leave inventory after equipping.`)
	}

	slot := equip.Equip.Slot
	if equipped := equipper.Equipper.Body().Slots[slot]; equipped != equip {
		t.Errorf(`Equipped item was %v, want %v`, equipped, equip)
	}
}

func TestEquipIntoOccupiedSlot(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip1 := g.NewObj(atItemSpec)
	equip2 := g.NewObj(atItemSpec)

	inv := equipper.Packer.Inventory()
	inv.Add(equip1)
	inv.Add(equip2)

	// Wield equip1
	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)
	// Wield equip2, swapping out equip1
	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(0)

	if swapped := inv.Top(); swapped != equip1 {
		t.Errorf(`First wield was not swapped out; got %v, want %v.`, swapped, equip1)
	}

	slot := equip2.Equip.Slot
	if equipped := equipper.Equipper.Body().Slots[slot]; equipped != equip2 {
		t.Errorf(`Equipped item was %v, want %v`, equipped, equip2)
	}
}

func TestEquipOutOfBounds(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryEquip()

	equip := g.NewObj(atItemSpec)
	inv := equipper.Packer.Inventory()
	inv.Add(equip)

	equipper.Equipper.TryEquip()
	equipper.Equipper.Equip(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Was mode %v after equip; want %v`, mode, ModeHud)
	}
}

func TestTryRemoveNothingEquipped(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeHud)
	}
}

func TestTryRemoveSomethingEquipped(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeRemove {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeRemove)
	}
}

func TestRemove(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(equip.Equip.Slot)

	if removed := equipper.Equipper.Body().Slots[equip.Equip.Slot]; removed != nil {
		t.Errorf(`Found %v in removed slot; want nil`, removed)
	}

	if removed := equipper.Packer.Inventory().Top(); removed != equip {
		t.Errorf(`Found %v in pack; want %v`, removed, equip)
	}
}

func TestRemoveOutOfBounds(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)

	equipper.Equipper.Body().Wear(equip)
	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(100)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Out-of-bounds removed switched mode to %v; want %v`, mode, ModeHud)
	}
}

func TestRemoveOverflowsToGround(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)
	equipper.Equipper.Body().Wear(equip)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(equipper, math.Pt(0, 0))

	equipper.Packer.Inventory().capacity = 0

	equipper.Equipper.TryRemove()
	equipper.Equipper.Remove(equip.Equip.Slot)

	if removed := equipper.Packer.Inventory().Top(); removed != nil {
		t.Errorf(`Found %v in pack; want nil`, removed)
	}

	if removed := equipper.Tile.Items.Top(); removed != equip {
		t.Errorf(`Found %v on floor; want %v`, removed, equip)
	}
}

func TestRemoveWhenInvAndGroundAreFull(t *testing.T) {
	g := NewGame()

	equipper := g.NewObj(atActorSpec)
	equip := g.NewObj(atItemSpec)
	equipper.Equipper.Body().Wear(equip)

	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(equipper, math.Pt(0, 0))

	equipper.Packer.Inventory().capacity = 0
	equipper.Tile.Items.capacity = 0

	equipper.Equipper.TryRemove()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryRemove switched to mode %v; want %v`, mode, ModeHud)
	}

}

func TestTryDropWithNothingInInventory(t *testing.T) {
	g := NewGame()

	packer := g.NewObj(atActorSpec)
	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryDrop w no items switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryDropWithFullGround(t *testing.T) {
	g := NewGame()

	packer := g.NewObj(atActorSpec)
	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(packer, math.Pt(0, 0))

	packer.Tile.Items.capacity = 0
	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`TryDrop w full ground switched to mode %v, want %v`, mode, ModeHud)
	}
}

func TestTryDrop(t *testing.T) {
	g := NewGame()

	packer := g.NewObj(atActorSpec)
	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(packer, math.Pt(0, 0))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()

	if mode := g.mode; mode != ModeDrop {
		t.Errorf(`TryDrop switched to mode %v, want %v`, mode, ModeDrop)
	}
}

func TestDrop(t *testing.T) {
	g := NewGame()

	packer := g.NewObj(atActorSpec)
	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(packer, math.Pt(0, 0))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()
	packer.Packer.Drop(0)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Dropping switched mode to %v, want %v`, mode, ModeHud)
	}

	if dropped := packer.Tile.Items.Top(); dropped != item {
		t.Errorf(`Dropped item was %v; want %v`, dropped, item)
	}
}

func TestDropOutOfBounds(t *testing.T) {
	g := NewGame()

	packer := g.NewObj(atActorSpec)
	l := NewLevel(4, 4, nil, IdentLevel)
	l.Place(packer, math.Pt(0, 0))

	item := g.NewObj(atItemSpec)
	packer.Packer.Inventory().Add(item)

	packer.Packer.TryDrop()
	packer.Packer.Drop(5)

	if mode := g.mode; mode != ModeHud {
		t.Errorf(`Dropping switched mode to %v, want %v`, mode, ModeHud)
	}
}
