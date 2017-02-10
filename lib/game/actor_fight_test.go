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
				Sheet: NewMonsterSheet(MonsterSheet{
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

type applyBrandTest struct {
	atk  Effects
	def  Effects
	verb string
	dice int
}

func TestApplyBrand(t *testing.T) {
	var (
		oldEffectVerbs = EffectVerbs
		oldResistMap   = ResistMap
		oldVulnMap     = VulnMap
		oldBrands      = Brands
	)
	const (
		fakeBrand1 Effect = NumEffects + iota
		fakeBrand2
		fakeVuln1
		fakeResist1
	)

	Brands = Effects{fakeBrand1, fakeBrand2}
	ResistMap = map[Effect]Effect{fakeBrand1: fakeResist1}
	VulnMap = map[Effect]Effect{fakeBrand1: fakeVuln1}
	EffectVerbs = map[Effect]string{fakeBrand1: "frobs", fakeBrand2: "norfs"}

	restoreEffectsDeps := func() {
		ResistMap = oldResistMap
		VulnMap = oldVulnMap
		EffectVerbs = oldEffectVerbs
		Brands = oldBrands
	}
	defer restoreEffectsDeps()

	tests := []applyBrandTest{
		{atk: Effects{fakeBrand1}, def: Effects{}, verb: "frobs", dice: 1},
		{atk: Effects{fakeBrand1}, def: Effects{fakeResist1}, verb: "hits", dice: 0},
		{atk: Effects{fakeBrand1}, def: Effects{fakeVuln1}, verb: "*frobs*", dice: 2},
		{atk: Effects{}, def: Effects{fakeVuln1}, verb: "hits", dice: 0},
		{atk: Effects{}, def: Effects{fakeResist1}, verb: "hits", dice: 0},
		{atk: Effects{fakeBrand1, fakeBrand2}, def: Effects{}, verb: "norfs", dice: 2},
		{atk: Effects{fakeBrand1, fakeBrand2}, def: Effects{fakeVuln1}, verb: "*frobs*", dice: 3},
		{atk: Effects{fakeBrand1, fakeBrand2}, def: Effects{fakeResist1}, verb: "norfs", dice: 1},
	}

	// No resist.
	for i, test := range tests {
		if dice, verb := applybrands(test.atk, test.def); dice != test.dice || verb != test.verb {
			t.Errorf(`Test %d: got (%d,"%s") want (%d,"%s")`, i, dice, verb, test.dice, test.verb)
		}
	}
}
