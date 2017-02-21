package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

type fakefighter struct {
	Trait
	Called bool
}

func (f *fakefighter) Hit(other Fighter) {
	f.Called = true
}

func TestPlayerMonsterCollisionsHit(t *testing.T) {
	g := newTestGame()
	pf := &fakefighter{Trait: Trait{obj: g.Player}}
	g.Player.Fighter = pf

	monster := g.NewObj(atActorSpec)
	mf := &fakefighter{Trait: Trait{obj: monster}}
	monster.Fighter = mf

	g.Level.Place(g.Player, math.Pt(1, 1))
	g.Level.Place(monster, math.Pt(1, 2))

	g.Player.Mover.Move(math.Pt(0, 1))

	if !pf.Called {
		t.Error("Moving player into other did not try to hit.")
	}

	monster.Mover.Move(math.Pt(0, -1))

	if !mf.Called {
		t.Error("Moving other into player did not try to hit.")
	}
}

func TestMonsterMonsterCollisionsHit(t *testing.T) {
	g := newTestGame()
	mon1 := g.NewObj(atActorSpec)
	mf1 := &fakefighter{Trait: Trait{obj: mon1}}
	mon1.Fighter = mf1

	mon2 := g.NewObj(atActorSpec)
	mf2 := &fakefighter{Trait: Trait{obj: mon2}}
	mon2.Fighter = mf2

	g.Level.Place(mon1, math.Pt(1, 1))
	g.Level.Place(mon2, math.Pt(1, 2))

	mon1.Mover.Move(math.Pt(0, 1))

	if mf1.Called {
		t.Error("Moving monster into monster tried to hit.")
	}
}

type hitTest struct {
	rolls    []int
	wanthp   int
	protdice Dice
}

func TestHit(t *testing.T) {
	tests := []hitTest{
		// Meleeroll = 1, evaderoll = 1, miss
		{[]int{1, 1}, 20, ZeroDice},
		// Meleeroll = 2, evaderoll = 1, roll 5 damage
		{[]int{2, 1, 5}, 15, ZeroDice},
		// Meleeroll = 2, evaderoll = 1, roll 5 damage, roll 2 prot
		{[]int{2, 1, 5, 2}, 17, NewDice(1, 4)},
		// Meleeroll = 8, evaderoll = 1, crit = 1, roll 5 + 3 damage
		{[]int{12, 1, 5, 3}, 12, ZeroDice},
		// Meleeroll = 15, evaderoll = 1, crit = 2, roll 3 + 2 + 1 damage
		{[]int{15, 1, 3, 2, 1}, 14, ZeroDice},
	}

	for i, test := range tests {
		testMonSpec := &Spec{
			Family:  FamActor,
			Genus:   GenMonster,
			Species: SpecOrc,
			Name:    "ORC",
			Traits: &Traits{
				Fighter: NewActorFighter,
				Sheet: NewMonsterSheet(&MonsterSheet{
					critdivmod: 0,
					maxhp:      20,
					damroll:    NewDice(1, 5),
					protroll:   test.protdice,
				}),
			},
		}

		g := newTestGame()
		attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
		FixRandomDie(test.rolls)
		defer RestoreRandom()

		attacker.Fighter.Hit(defender.Fighter)
		if hp := defender.Sheet.HP(); hp != test.wanthp {
			t.Errorf(`Test %d: Defender has %d hp; want %d.`, i, hp, test.wanthp)
		}
	}
}

func TestHitPoison(t *testing.T) {
	testMonSpec := &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits: &Traits{
			Fighter: NewActorFighter,
			Sheet: NewMonsterSheet(&MonsterSheet{
				critdivmod: 0,
				maxhp:      20,
				damroll:    NewDice(1, 5),
				atkeffects: NewEffects(map[Effect]int{BrandPoison: 1}),
			}),
			Ticker: NewActorTicker,
		},
	}

	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, 2 of which is poison.
	FixRandomDie([]int{7, 1, 5, 2})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if poison := defender.Ticker.Counter(EffectPoison); poison != 2 {
		t.Errorf(`Poison counter %d; want 2`, poison)
	}
}

func TestHitStun(t *testing.T) {
	testMonSpec := &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits: &Traits{
			Fighter: NewActorFighter,
			Sheet: NewMonsterSheet(&MonsterSheet{
				critdivmod: 0,
				maxhp:      20,
				damroll:    NewDice(1, 5),
				atkeffects: NewEffects(map[Effect]int{EffectStun: 1}),
			}),
			Ticker: NewActorTicker,
		},
	}

	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, and make sure to win stun skillroll (10 vs 0 on d10s.)
	FixRandomDie([]int{7, 1, 5, 10, 0})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if stun := defender.Ticker.Counter(EffectStun); stun != 5 {
		t.Errorf(`Stun counter %d; want 5`, stun)
	}
}

//type applyBrandTest struct {
//	atk  Effects
//	def  Effects
//	verb string
//	dice int
//}
//
//func TestApplyBrand(t *testing.T) {
//	const (
//		fakeBrand1 Effect = NumEffects + iota
//		fakeBrand2
//		fakeResist1
//	)
//
//	oldEffectsSpecs := EffectsSpecs
//
//	EffectsSpecs = EffectsSpec{
//		fakeBrand1:  {Type: EffectTypeBrand, ResistedBy: fakeResist1, Verb: "frobs"},
//		fakeBrand2:  {Type: EffectTypeBrand, Verb: "norfs"},
//		fakeResist1: {Type: EffectTypeResist},
//	}
//
//	restoreEffectsDeps := func() {
//		EffectsSpecs = oldEffectsSpecs
//	}
//	defer restoreEffectsDeps()
//
//	tests := []applyBrandTest{
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1}),
//			def:  Effects{},
//			verb: "frobs",
//			dice: 1,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 2}),
//			def:  Effects{},
//			verb: "frobs",
//			dice: 1,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: 1}),
//			verb: "hits",
//			dice: 0,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: 2}),
//			verb: "hits",
//			dice: 0,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: -1}),
//			verb: "*frobs*",
//			dice: 2,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: -2}),
//			verb: "*frobs*",
//			dice: 2,
//		},
//		{
//			atk:  Effects{},
//			def:  NewEffects(map[Effect]int{fakeResist1: -1}),
//			verb: "hits",
//			dice: 0,
//		},
//		{
//			atk:  Effects{},
//			def:  NewEffects(map[Effect]int{fakeResist1: 1}),
//			verb: "hits",
//			dice: 0,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1, fakeBrand2: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: -1}),
//			verb: "*frobs*",
//			dice: 3,
//		},
//		{
//			atk:  NewEffects(map[Effect]int{fakeBrand1: 1, fakeBrand2: 1}),
//			def:  NewEffects(map[Effect]int{fakeResist1: 1}),
//			verb: "norfs",
//			dice: 1,
//		},
//	}
//
//	// No resist.
//	for i, test := range tests {
//		if dice, verb := applybrands(test.atk, test.def); dice != test.dice || verb != test.verb {
//			t.Errorf(`Test %d: got (%d,"%s") want (%d,"%s")`, i, dice, verb, test.dice, test.verb)
//		}
//	}
//}
