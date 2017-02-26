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

	atkroll := combatroll(attacker.Obj()) + atk.Melee
	defroll := combatroll(defender.Obj()) + def.Evasion
	residual := atkroll - defroll

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

	if dmg <= 0 {
		return
	}

	for effect, _ := range atk.Effects {
		switch effect {
		case BrandPoison:
			defender.Obj().Ticker.AddEffect(EffectPoison, poisondmg)
		case EffectStun:
			score := atk.CritDiv - BaseCritDiv + attacker.Obj().Sheet.Stat(Str)
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := def.Effects.Resists(effect)
			won, _ := skillcheck(score, difficulty, resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectStun, dmg)
			}
		case EffectBlind:
			score := 10
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := resistmod(def.Effects.Resists(effect))
			won, _ := skillcheck(score, difficulty, resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectBlind, DieRoll(5, 4))
			}
		case EffectConfuse:
			// TODO: Eventually remove this check and instead use a Cruel-Blow
			// style check, cruel blow should be the only ability that gives
			// confusion melee anyways.
			score := 10
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := resistmod(def.Effects.Resists(effect))
			won, _ := skillcheck(score, difficulty, resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectConfuse, DieRoll(5, 4))
			}
		case EffectPara:
			// Don't let the effect accumulate; also, give the defender a
			// chance to break out when they are hurt.
			if defender.Obj().Sheet.Paralyzed() {
				checkpara(defender)
				break
			}

			score := 10
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := resistmod(def.Effects.Resists(effect))
			won, _ := skillcheck(score, difficulty, resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectPara, DieRoll(4, 4))
			}
		case EffectPetrify:
			// Don't let the effect accumulate.
			if defender.Obj().Sheet.Petrified() {
				break
			}
			score := 10
			difficulty := defender.Obj().Sheet.Skill(Chi)
			resists := resistmod(def.Effects.Resists(effect))
			won, _ := skillcheck(score, difficulty, resists, attacker.Obj(), defender.Obj())
			if won {
				defender.Obj().Ticker.AddEffect(EffectPetrify, DieRoll(4, 4))
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
	if basedmg == 0 {
		return 0, 0, "hits"
	}

	brands := atk.Brands()
	fixverb := false
	branddmg, poisondmg, verb = 0, 0, "hits"

	for brand, info := range brands {
		raw := DieRoll(1, basedmg)
		resisted := def.ResistDmg(brand, raw)

		if brand == BrandPoison {
			poisondmg += resisted
		} else {
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
	return branddmg, poisondmg, verb
}

func checkpara(defender Fighter) {
	obj := defender.Obj()
	sheet := obj.Sheet
	if !sheet.Paralyzed() {
		return
	}

	if OneIn(2) {
		msg := fmt.Sprintf("%s breaks out of paralysis!", obj.Spec.Name)
		obj.Game.Events.Message(msg)

		sheet.SetParalyzed(false)
		obj.Ticker.RemoveEffect(EffectPara)
	}
}
