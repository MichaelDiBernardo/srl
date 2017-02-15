package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
}

// Melee combat.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *ActorFighter) Hit(other Fighter) {
	hit(f, other)
}

func hit(attacker Fighter, defender Fighter) {
	atk, def := attacker.Obj().Sheet.Attack(), defender.Obj().Sheet.Defense()
	residual := DieRoll(1, 20) + atk.Melee - DieRoll(1, 20) + def.Evasion
	aname, dname := attacker.Obj().Spec.Name, defender.Obj().Spec.Name

	if residual <= 0 {
		msg := fmt.Sprintf("%v missed %v.", aname, dname)
		attacker.Obj().Game.Events.Message(msg)
		return
	}

	crits := residual / atk.CritDiv
	extra, verb := applybrands(atk.Effects, def.Effects)
	dmg := math.Max(0, atk.RollDamage(crits+extra)-def.RollProt())

	critstr := ""
	if crits > 0 {
		critstr = fmt.Sprintf(" %dx critical!", crits)
	}
	msg := fmt.Sprintf("%s %s %s (%d).%s", aname, verb, dname, dmg, critstr)
	attacker.Obj().Game.Events.Message(msg)

	for effect, _ := range atk.Effects {
		switch effect {
		case BrandPoison:
			defender.Obj().Ticker.AddEffect(EffectPoison, dmg)
			dmg = 0
		}
	}
	defender.Obj().Sheet.Hurt(dmg)
}

func applybrands(atk Effects, def Effects) (dice int, verb string) {
	brands := atk.Brands()
	fixverb := false
	dice, verb = 0, "hits"

	for brand, info := range brands {
		// Extra dice from a brand:
		// Any amount of resist = 0
		// No resist = 1
		// Any amount of vuln = 2
		resists := def.Resists(brand)
		extradice := -(math.Sgn(resists) - 1)
		dice += extradice

		if extradice == 0 {
			continue
		}

		newverb := info.Verb

		// Vulnerability
		if extradice == 2 {
			fixverb = true
			verb = fmt.Sprintf("*%s*", newverb)
		}

		// If we've ever found a vulnerability before (including in this
		// iteration), keep the old verb because it's at least as important as
		// any other verb we might use. e.g. if you have fire and cold brands,
		// and the target is vuln to fire, we want '*burns*' to have priority
		// over 'freezes'
		if !fixverb {
			verb = newverb
		}
	}
	return dice, verb
}
