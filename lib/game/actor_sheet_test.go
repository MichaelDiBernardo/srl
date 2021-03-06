package game

import (
	"testing"
)

func TestPlayerMaxHPCalc(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 1}},
	})

	if maxhp, want := obj.Sheet.MaxHP(), 20; maxhp != want {
		t.Errorf(`MaxHP() was %d, want %d`, maxhp, want)
	}
}

func TestPlayerMaxMPCalc(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Mnd: 2}},
	})

	if maxmp, want := obj.Sheet.MaxMP(), 30; maxmp != want {
		t.Errorf(`MaxMP() was %d, want %d`, maxmp, want)
	}
}

func TestHurtingPlayerToDeathEndsGame(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})
	g.Player = obj

	obj.Sheet.Hurt(9999999)
	if m := g.mode; m != ModeGameOver {
		t.Errorf(`Killing player changed mode to %v; want %v`, m, ModeGameOver)
	}
}

func TestHealPlayerDoesntExceedMaxHP(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}, skills: &skills{}, hp: 1})
	g.Player = obj

	obj.Sheet.Heal(9999999)
	if hp, maxhp := obj.Sheet.HP(), obj.Sheet.MaxHP(); hp != maxhp {
		t.Errorf(`Player healed to %d hp; want %d.`, hp, maxhp)
	}
}

func TestChangingPlayerVitAdjustsHP(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Vit: 2}},
		hp:    15,
	})

	g.Player = obj

	obj.Sheet.ChangeStatMod(Vit, 3)
	if hp := obj.Sheet.HP(); hp != 30 {
		t.Errorf(`Player hp scaled to %d hp; want %d.`, hp, 30)
	}

	obj.Sheet.SetStat(Vit, -1)
	if hp := obj.Sheet.HP(); hp != 15 {
		t.Errorf(`Player hp scaled to %d hp; want %d.`, hp, 15)
	}
}

func TestChangingPlayerMndAdjustsMP(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Mnd: 2}},
		mp:    15,
	})

	g.Player = obj

	obj.Sheet.ChangeStatMod(Mnd, 3)
	if mp := obj.Sheet.MP(); mp != 30 {
		t.Errorf(`Player mp scaled to %d mp; want %d.`, mp, 30)
	}

	obj.Sheet.SetStat(Mnd, -1)
	if mp := obj.Sheet.MP(); mp != 15 {
		t.Errorf(`Player mp scaled to %d mp; want %d.`, mp, 15)
	}
}

var astKnifeSpec = &Spec{
	Family:  FamItem,
	Genus:   GenEquipment,
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
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 7)})
}

func TestPlayerAttackStrBonusBelowCap(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Str: 1}},
	})

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 8)})
}

func TestPlayerAttackStrBonusAboveCap(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Str: 3}},
	})

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 1, Damroll: NewDice(1, 9)})
}

func TestPlayerAttackMeleeBonus(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Agi: 3}},
	})

	weap := g.NewObj(astKnifeSpec)
	obj.Equipper.Body().Wear(weap)

	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 4, Damroll: NewDice(1, 7)})
}

func TestPlayerAttackFistNoStrSides(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})
	atk := obj.Sheet.Attack()
	testAtkEq(t, atk, Attack{Melee: 0, Damroll: NewDice(1, 1)})
}

func TestPlayerAttackFistStrSides(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Str: 10}},
	})

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
	// This is good enough for the cases we have.
	for _, wd := range want.ProtDice {
		found := false
		for _, ad := range def.ProtDice {
			if ad == wd {
				found = true
			}
		}
		if !found {
			t.Errorf(`def.ProtDice was %+v, want %+v`, def.ProtDice, want.ProtDice)
		}
	}
	if l, w := len(def.CorrDice), len(want.CorrDice); l != w {
		t.Errorf(`len(def.CorrDice) was %d, want %d`, l, w)
	}
	// This is good enough for the cases we have.
	for _, wd := range want.CorrDice {
		found := false
		for _, ad := range def.CorrDice {
			if ad == wd {
				found = true
			}
		}
		if !found {
			t.Errorf(`def.CorrDice was %+v, want %+v`, def.CorrDice, want.CorrDice)
		}
	}
}

func TestPlayerDefenseNoArmorOrEvasion(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{Evasion: 0})
}

func TestPlayerDefenseNoArmorWithEvasion(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		stats: &stats{stats: statlist{Agi: 2}},
	})

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{Evasion: 2})
}

func TestPlayerDefenseWithArmor(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})

	armspec1 := &Spec{
		Family:  FamItem,
		Genus:   GenEquipment,
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
		Genus:   GenEquipment,
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

func TestPlayerDefenseCorr(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	obj.Sheet = NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}})
	obj.Sheet.SetCorrosion(2)

	def := obj.Sheet.Defense()
	testDefEq(t, def, Defense{Evasion: 0, CorrDice: []Dice{NewDice(2, 4)}})
}

func TestPlayerBlind(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	sheet := NewPlayerSheetFromSpec(&PlayerSheet{
		Trait:  Trait{obj: obj},
		blind:  true,
		skills: &skills{skills: skilllist{Melee: 10, Evasion: 10, Shooting: 10}},
	})
	obj.Sheet = sheet

	for _, skill := range []SkillName{Melee, Evasion, Shooting} {
		if sk := sheet.Skill(skill); sk != 5 {
			t.Errorf(`Skill %v was %d, want 5`, skill, sk)
		}
	}

	if s := sheet.Sight(); s != 0 {
		t.Errorf(`sheet.Sight() was %d, want 0`, s)
	}

	if m := sheet.Attack().Melee; m != 5 {
		t.Errorf(`atk.Melee = %d, want 5`, m)
	}

	if e := sheet.Defense().Evasion; e != 5 {
		t.Errorf(`def.Evasion = %d, want 5`, e)
	}
}

func TestPlayerSlow(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	sheet := NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
		slow:  true,
		speed: 2,
	})
	obj.Sheet = sheet

	if s := sheet.Speed(); s != 1 {
		t.Errorf(`sheet.Speed() was %d, want 1`, s)
	}

	sheet.speed = 1

	if s := sheet.Speed(); s != 1 {
		t.Errorf(`sheet.Speed() was %d, want 1`, s)
	}
}

func TestMonsterSlow(t *testing.T) {
	g := newTestGame()
	obj := g.NewObj(&Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits:  &Traits{Sheet: NewMonsterSheet(&MonsterSheet{speed: 2, slow: true})},
	})
	sheet := obj.Sheet.(*MonsterSheet)

	if s := sheet.Speed(); s != 1 {
		t.Errorf(`sheet.Speed() was %d, want 1`, s)
	}

	sheet.speed = 1

	if s := sheet.Speed(); s != 1 {
		t.Errorf(`sheet.Speed() was %d, want 1`, s)
	}
}

func TestPlayerConfuse(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	sheet := NewPlayerSheetFromSpec(&PlayerSheet{
		Trait: Trait{obj: obj},
	})
	obj.Sheet = sheet

	sheet.SetConfused(true)
	for i := Chi; i < NumSkills; i++ {
		if mod := sheet.SkillMod(i); mod != -5 {
			t.Errorf(`Skillmod %v was %d, want -5`, i, mod)
		}
	}

	sheet.SetConfused(false)
	for i := Chi; i < NumSkills; i++ {
		if mod := sheet.SkillMod(i); mod != 0 {
			t.Errorf(`Skillmod %v was %d, want 0`, i, mod)
		}
	}
}

func TestPlayerPara(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	sheet := NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}, hp: 10})
	obj.Sheet = sheet

	sheet.SetParalyzed(true)

	if ev := sheet.Skill(Evasion); ev != dumpsterEvasion {
		t.Errorf(`Para evasion %d, want %d`, ev, dumpsterEvasion)
	}

	if ev := sheet.Defense().Evasion; ev != dumpsterEvasion {
		t.Errorf(`Para Defense().Evasion %d, want %d`, ev, dumpsterEvasion)
	}

	if sheet.CanAct() {
		t.Error(`Para CanAct() was true, want false`)
	}
}

func TestPlayerPetrify(t *testing.T) {
	g := newTestGame()
	obj := g.Player
	sheet := NewPlayerSheetFromSpec(&PlayerSheet{Trait: Trait{obj: obj}, hp: 10})
	obj.Sheet = sheet

	sheet.SetPetrified(true)

	if ev := sheet.Skill(Evasion); ev != dumpsterEvasion {
		t.Errorf(`Petrify evasion %d, want %d`, ev, dumpsterEvasion)
	}

	if ev := sheet.Defense().Evasion; ev != dumpsterEvasion {
		t.Errorf(`Petrify Defense().Evasion %d, want %d`, ev, dumpsterEvasion)
	}

	if prot := sheet.Defense().ProtDice; len(prot) != 1 && prot[0] != petrifyProt {
		t.Errorf(`Petrify Defense().ProtDice was %v, want %v`, prot, []Dice{petrifyProt})
	}
	if cr := sheet.Defense().Effects.Has(ResistCrit); cr != petrifyCritResist {
		t.Errorf(`Petrify Defense().Effects[ResistCrit] was %d, want %d`, cr, petrifyCritResist)
	}

	if regen := sheet.Regen(); regen != 0 {
		t.Errorf(`Petrify Regen() was %d, want 0`, regen)
	}

	if sheet.CanAct() {
		t.Error(`Petrify CanAct() was true, want false`)
	}
}
