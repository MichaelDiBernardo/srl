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

func TestCantHitWhenAfraid(t *testing.T) {
	g := newTestGame()
	monster := g.NewObj(atActorSpec)

	g.Player.Sheet.SetAfraid(true)

	g.Level.Place(g.Player, math.Pt(1, 1))
	g.Level.Place(monster, math.Pt(1, 2))

	ok, err := g.Player.Mover.Move(math.Pt(0, 1))

	if ok {
		t.Error(`Attacking while afraid used a turn.`)
	}

	if err != ErrTooScaredToHit {
		t.Errorf(`Attacking while afraid gave error %v, want %v`, err, ErrTooScaredToHit)
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

func makeTestHitterSpec(atk Effects) *Spec {
	return &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: SpecOrc,
		Name:    "ORC",
		Traits: &Traits{
			Fighter: NewActorFighter,
			Sheet: NewMonsterSheet(&MonsterSheet{
				critdivmod: 0,
				maxhp:      20,
				maxmp:      10,
				damroll:    NewDice(1, 5),
				speed:      1,
				atkeffects: atk,
			}),
			Ticker: NewActorTicker,
		},
	}
}

func TestHitCritResist(t *testing.T) {
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
				defeffects: NewEffects(map[Effect]int{ResistCrit: 1}),
			}),
		},
	}

	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)

	// Roll a residual of 7, which should normally result in 1 crit. but in In
	// this case, though, the divisor should be 8 (7 + 1 crit resist). This
	// should result in a 1d5 damroll, which we fix at 3 dmg.
	FixRandomDie([]int{10, 3, 3})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)
	if hp, want := defender.Sheet.HP(), 17; hp != want {
		t.Errorf(`Defender has %d hp; want %d.`, hp, want)
	}
}

func TestHitVamp(t *testing.T) {
	g := newTestGame()
	spec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectVamp: 1}))
	attacker, defender := g.NewObj(spec), g.NewObj(spec)
	attacker.Sheet.setHP(1)
	attacker.Sheet.setMP(1)

	// We're gonna kill an actor here (actors were harmed in the making of this
	// test), so let's place them to keep stuff from crashing when we try to
	// remove the dead one on death.
	g.Level.Place(attacker, math.Pt(1, 1))
	g.Level.Place(defender, math.Pt(2, 1))

	// Hit for a bazillion damage so we can guarantee the target is dead.
	FixRandomDie([]int{4, 3, 300})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)
	if hp, want := attacker.Sheet.HP(), 6; hp != want {
		t.Errorf(`Attacker has %d hp; want %d.`, hp, want)
	}

	if mp, want := attacker.Sheet.MP(), 3; mp != want {
		t.Errorf(`Attacker has %d mp; want %d.`, mp, want)
	}
}

func TestHitPoison(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{BrandPoison: 1}))
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

func TestHitAcid(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{BrandAcid: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, then 1 for extra brand damage, then 1 for the OneIn
	// check, then 5 turns of corr.
	FixRandomDie([]int{7, 1, 5, 1, 1, 1, 1, 2, 1})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if corr := defender.Ticker.Counter(EffectShatter); corr != 5 {
		t.Errorf(`Shatter counter %d; want 5`, corr)
	}
}

func TestHitStun(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectStun: 1}))
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

func TestHitShatter(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectShatter: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	attacker.Sheet.(*MonsterSheet).critdivmod = 10
	// Roll 5 damage, and make sure to win shatter skillroll (10 vs 0 on d10s.)
	FixRandomDie([]int{7, 1, 5, 10, 0, 1, 1, 1, 2})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if corr := defender.Ticker.Counter(EffectShatter); corr != 5 {
		t.Errorf(`Shatter counter %d; want 5`, corr)
	}
}

func TestHitCut(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectCut: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 2 crits, do 4 damage (2 + 1 + 1), roll 1 on the crit check to force
	// cut.
	FixRandomDie([]int{20, 1, 2, 1, 1, 1})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if cut := defender.Ticker.Counter(EffectCut); cut != 2 {
		t.Errorf(`Cut counter %d; want 2`, cut)
	}
}

func TestHitPara(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectPara: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, and make sure to win skillroll (10 vs 0 on d10s. We then
	// roll 12 para turns (3 * 4).)
	// We'll hit a second time, this time hitting for 1 damage. We fix the die
	// so that it breaks the defender out of para.
	// We use 2 as our second attack roll because the defender will already be
	// at -5 due to being paralyzed.
	FixRandomDie([]int{7, 1, 5, 10, 0, 3, 3, 3, 3 /* second attack */, 2, 1, 1, 1})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if para := defender.Ticker.Counter(EffectPara); para != 12 {
		t.Errorf(`Para counter %d; want 12`, para)
	}

	attacker.Fighter.Hit(defender.Fighter)

	if para := defender.Ticker.Counter(EffectPara); para != 0 {
		t.Errorf(`Para counter %d; want 0`, para)
	}
	if defender.Sheet.Paralyzed() {
		t.Error(`defender.Sheet.Paralyzed() was true, want false`)
	}
}

func TestHitBlind(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectBlind: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, and make sure to win skillroll (10 vs 0 on d10s. We then
	// roll 15 blind turns (3 * 5).)
	FixRandomDie([]int{7, 1, 5, 10, 0, 3, 3, 3, 3, 3})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if blind := defender.Ticker.Counter(EffectBlind); blind != 15 {
		t.Errorf(`Blind counter %d; want 15`, blind)
	}
}

func TestHitConfuse(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectConfuse: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, and make sure to win skillroll (10 vs 0 on d10s. We then
	// roll 15 confusion turns (3 * 5).)
	FixRandomDie([]int{7, 1, 5, 10, 0, 3, 3, 3, 3, 3})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if conf := defender.Ticker.Counter(EffectConfuse); conf != 15 {
		t.Errorf(`Conf counter %d; want 15`, conf)
	}
}

func TestHitPetrify(t *testing.T) {
	testMonSpec := makeTestHitterSpec(NewEffects(map[Effect]int{EffectPetrify: 1}))
	g := newTestGame()
	attacker, defender := g.NewObj(testMonSpec), g.NewObj(testMonSpec)
	// Roll 5 damage, and make sure to win skillroll (10 vs 0 on d10s. We then
	// roll 12 petrify turns (3 * 4).)
	FixRandomDie([]int{7, 1, 5, 10, 0, 3, 3, 3, 3})
	defer RestoreRandom()

	attacker.Fighter.Hit(defender.Fighter)

	if petr := defender.Ticker.Counter(EffectPetrify); petr != 12 {
		t.Errorf(`Petrify counter %d; want 12`, petr)
	}
}

// Test raw damage due to brands.
type applyBrandTest struct {
	// Inputs
	atk   Effects
	def   Effects
	rolls []int
	// Expected
	branddmg int
	verb     string
}

func TestApplyBrand(t *testing.T) {
	const (
		fakeBrand1 Effect = NumEffects + iota
		fakeBrand2
		fakeResist1
	)

	oldEffectsSpecs := EffectsSpecs

	EffectsSpecs = EffectsSpec{
		fakeBrand1:  {Type: EffectTypeBrand, ResistedBy: fakeResist1, Verb: "frobs"},
		fakeBrand2:  {Type: EffectTypeBrand, Verb: "norfs"},
		fakeResist1: {Type: EffectTypeResist},
	}

	restoreEffectsDeps := func() {
		EffectsSpecs = oldEffectsSpecs
	}
	defer restoreEffectsDeps()

	tests := []applyBrandTest{
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1}),
			def:      Effects{},
			rolls:    []int{5},
			verb:     "frobs",
			branddmg: 5,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 2}),
			def:      Effects{},
			rolls:    []int{5},
			verb:     "frobs",
			branddmg: 5,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: 1}),
			rolls:    []int{5},
			verb:     "frobs",
			branddmg: 2,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: 2}),
			rolls:    []int{5},
			verb:     "frobs",
			branddmg: 1,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: -1}),
			rolls:    []int{5},
			verb:     "*frobs*",
			branddmg: 10,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: -2}),
			rolls:    []int{5},
			verb:     "*frobs*",
			branddmg: 15,
		},
		{
			atk:      NewEffects(map[Effect]int{}),
			def:      NewEffects(map[Effect]int{fakeResist1: -1}),
			rolls:    []int{5},
			verb:     "hits",
			branddmg: 0,
		},
		{
			atk:      NewEffects(map[Effect]int{}),
			def:      NewEffects(map[Effect]int{fakeResist1: 1}),
			rolls:    []int{5},
			verb:     "hits",
			branddmg: 0,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1, fakeBrand2: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: 1}),
			rolls:    []int{5, 5},
			branddmg: 7,
		},
		{
			atk:      NewEffects(map[Effect]int{fakeBrand1: 1, fakeBrand2: 1}),
			def:      NewEffects(map[Effect]int{fakeResist1: -1}),
			rolls:    []int{5, 5},
			verb:     "*frobs*",
			branddmg: 15,
		},
	}

	// No resist.
	for i, test := range tests {
		func() {
			FixRandomDie(test.rolls)
			defer RestoreRandom()
			if branddmg, _, verb := applybs(10, test.atk, test.def); branddmg != test.branddmg || (test.verb != "" && verb != test.verb) {
				t.Errorf(`Test %d: got (%d, "%s") want (%d, "%s")`, i, branddmg, verb, test.branddmg, test.verb)
			}
		}()
	}
}

func TestApplyBrandPoison(t *testing.T) {
	FixRandomDie([]int{5, 5, 5, 5})
	defer RestoreRandom()

	if branddmg, poisondmg, _ := applybs(10, NewEffects(map[Effect]int{BrandPoison: 1}), NewEffects(map[Effect]int{})); branddmg != 0 || poisondmg != 5 {
		t.Errorf(`applybrand poisondmg: got (%d, %d) want (0, 5)`, branddmg, poisondmg)
	}
	if branddmg, poisondmg, _ := applybs(10, NewEffects(map[Effect]int{BrandPoison: 1}), NewEffects(map[Effect]int{ResistPoison: 1})); branddmg != 0 || poisondmg != 2 {
		t.Errorf(`applybrand poisondmg: got (%d, %d) want (0, 2)`, branddmg, poisondmg)
	}
	if branddmg, poisondmg, _ := applybs(10, NewEffects(map[Effect]int{BrandPoison: 1, BrandFire: 1}), NewEffects(map[Effect]int{})); branddmg != 5 || poisondmg != 5 {
		t.Errorf(`applybrand poisondmg: got (%d, %d) want (5, 5)`, branddmg, poisondmg)
	}
}
