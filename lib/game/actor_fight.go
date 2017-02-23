package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
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
	// Calculate raw phys damage.
	dmg := math.Max(0, atk.RollDamage(crits)-def.RollProt())
	// Figure out how much branded damage we did.
	branddmg, poisondmg, verb := applybrands(dmg, atk.Effects, def.Effects)
	dmg += branddmg

	critstr := ""
	if crits > 0 {
		critstr = fmt.Sprintf(" %dx critical!", crits)
	}
	msg := fmt.Sprintf("%s %s %s (%d).%s", aname, verb, dname, dmg, critstr)
	attacker.Obj().Game.Events.Message(msg)

	for effect, _ := range atk.Effects {
		switch effect {
		case BrandPoison:
			defender.Obj().Ticker.AddEffect(EffectPoison, poisondmg)
		case EffectStun:
			score := atk.CritDiv - BaseCritDiv + attacker.Obj().Sheet.Stat(Str)
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := resistmod(def.Effects.Resists(effect))
			won, _ := skillcheck(score, difficulty+resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectStun, dmg)
			}
		case EffectCut:
			if crits > DieRoll(1, 2) {
				defender.Obj().Ticker.AddEffect(EffectCut, dmg/2)
			}
		}
	}
	defender.Obj().Sheet.Hurt(dmg)
}

// Given the base physical damage done by an attack, and the atk and def
// effects, this figures out how much raw and poison damage should be done from
// brands. Poison damage is separated out because it is applied as
// damage-over-time, instead of being immediately inflicted on the target.
func applybrands(basedmg int, atk Effects, def Effects) (branddmg, poisondmg int, verb string) {
	brands := atk.Brands()
	fixverb := false
	branddmg, poisondmg, verb = 0, 0, "hits"
	log.Printf("applybrand: %d %v vs %v", basedmg, atk, def)

	for brand, info := range brands {
		raw := DieRoll(1, basedmg)
		resisted := def.ResistDmg(brand, raw)

		if brand == BrandPoison {
			poisondmg += resisted
			log.Printf("\tDid raw:%d resisted:%d poison", raw, resisted)
		} else {
			log.Printf("\tDid raw:%d resisted:%d for brand %v", raw, resisted, brand)
			branddmg += resisted
		}

		newverb := info.Verb

		// Vulnerability
		if def.Resists(brand) < 0 {
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
	log.Printf("\tReturning branddmg:%d poisondmg:%d verb:%s", branddmg, poisondmg, verb)
	return branddmg, poisondmg, verb
}
