package game

import (
	"testing"
)

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

var astKnifeSpec = &Spec{
	Family:  FamItem,
	Genus:   GenEquip,
	Species: SpecSword,
	Name:    "KNIFE",
	Traits: &Traits{
		Equipment: NewEquipment(Equipment{
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

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 7)})
}

func TestPlayerAttackStrBonusBelowCap(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 1, agi: 0}

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 8)})
}

func TestPlayerAttackStrBonusAboveCap(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 3, agi: 0}

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 9)})
}

func TestPlayerAttackMeleeBonus(t *testing.T) {
	g := NewGame()
	obj := g.NewObj(PlayerSpec)
	obj.Sheet = &PlayerSheet{Trait: Trait{obj: obj}, str: 0, agi: 3}

	weap := g.NewObj(astKnifeSpec)
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
			Equipment: NewEquipment(Equipment{
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
			Equipment: NewEquipment(Equipment{
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
