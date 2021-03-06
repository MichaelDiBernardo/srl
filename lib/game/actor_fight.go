package game

import (
	"fmt"
	"log"

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
	a, d := attacker.Obj(), defender.Obj()
	atk, def := a.Sheet.Attack(), d.Sheet.Defense()

	atkroll := combatroll(attacker.Obj()) + atk.Melee
	defroll := combatroll(defender.Obj()) + def.Evasion
	residual := atkroll - defroll

	aname, dname := a.Spec.Name, d.Spec.Name

	if residual <= 0 {
		msg := fmt.Sprintf("%v missed %v.", aname, dname)
		a.Game.Events.Message(msg)
		return
	}

	crits := residual / (atk.CritDiv + def.Effects.Has(ResistCrit))

	// Calculate raw phys damage.
	droll, proll := atk.RollDamage(crits), def.RollProt()
	log.Printf("DR: %d PR: %d", droll, proll)
	dmg := math.Max(0, droll-proll)

	// Figure out how much branded damage we did.
	xdmg, poisondmg := applybs(dmg, atk.Effects, def.Effects)
	dmg += xdmg

	critstr := ""
	if crits > 0 {
		critstr = fmt.Sprintf(" %dx critical!", crits)
	}

	msg := fmt.Sprintf("%s %s %s (%d).%s", aname, atk.Verb, dname, dmg, critstr)
	a.Game.Events.Message(msg)

	if dmg <= 0 {
		return
	}

	ispara := d.Sheet.Paralyzed()

	for effect, _ := range atk.Effects {
		switch effect {
		case BrandPoison:
			d.Ticker.AddEffect(EffectPoison, poisondmg)
		case BrandAcid:
			if OneIn(def.Effects.Resists(effect) + 1) {
				d.Ticker.AddEffect(EffectShatter, DieRoll(4, 4))
			}
		case EffectStun:
			score := atk.CritDiv - BaseCritDiv + a.Sheet.Stat(Str)
			difficulty := d.Sheet.Skill(Chi)
			resists := def.Effects.Resists(effect)
			won, _ := skillcheck(score, difficulty, resists, a, d)
			if won {
				d.Ticker.AddEffect(EffectStun, dmg)
			}
		case EffectBlind:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectBlind, DieRoll(5, 4))
			}
		case EffectConfuse:
			// TODO: Eventually remove this check and instead use a Cruel-Blow
			// style check, cruel blow should be the only ability that gives
			// confusion melee anyways.
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectConfuse, DieRoll(5, 4))
			}
		case EffectPara:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectPara, DieRoll(4, 4))
			}
		case EffectPetrify:
			if savingthrow(d, def.Effects, effect) {
				d.Ticker.AddEffect(EffectPetrify, DieRoll(4, 4))
			}
		case EffectCut:
			if crits > DieRoll(1, 2) {
				d.Ticker.AddEffect(EffectCut, dmg/2)
			}
		case EffectShatter:
			score := atk.CritDiv - BaseCritDiv + a.Sheet.Stat(Str)
			won, _ := skillcheck(score, 10, 0, a, d)
			if won {
				d.Ticker.AddEffect(EffectShatter, DieRoll(4, 4))
			}
		case EffectDrainStr:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainStr, 1)
			}
		case EffectDrainAgi:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainAgi, 1)
			}
		case EffectDrainVit:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainVit, 1)
			}
		case EffectDrainMnd:
			r := def.Effects.Resists(effect)
			won, _ := skillcheck(a.Sheet.Skill(Chi), d.Sheet.Skill(Chi), r, a, d)
			if won {
				d.Ticker.AddEffect(EffectDrainMnd, 1)
			}
		}
	}

	if ispara {
		checkpara(defender)
	}
	d.Sheet.Hurt(dmg)

	// Have to handle this outside the effect loop above because we need the
	// damage to be applied to the target first to determine if they're dead.
	if d.Sheet.Dead() &&
		atk.Effects.Has(EffectVamp) > 0 &&
		def.Effects.Resists(EffectVamp) <= 0 {
		vamp(a, d)
	}
}

// Given the base physical damage done by an attack, and the atk and def
// effects, this figures out how much extra and poison damage should be done from
// brands and slays. Poison damage is separated out because it is applied as
// damage-over-time, instead of being immediately inflicted on the target.
func applybs(basedmg int, atk Effects, def Effects) (xdmg, poisondmg int) {
	if basedmg == 0 {
		return 0, 0
	}

	slays, brands := atk.Slays(), atk.Brands()
	xdmg, poisondmg = 0, 0

	for slay, _ := range slays {
		if def.SlainBy(slay) <= 0 {
			return
		}
		xdmg += DieRoll(1, basedmg)
	}

	for brand, _ := range brands {
		raw := DieRoll(1, basedmg)
		resisted := def.ResistDmg(brand, raw)

		if brand == BrandPoison {
			poisondmg += resisted
		} else {
			xdmg += resisted
		}
	}
	return xdmg, poisondmg
}

func checkpara(defender Fighter) {
	obj := defender.Obj()
	sheet := obj.Sheet

	if OneIn(2) {
		msg := fmt.Sprintf("%s breaks out of paralysis!", obj.Spec.Name)
		obj.Game.Events.Message(msg)

		sheet.SetParalyzed(false)
		obj.Ticker.RemoveEffect(EffectPara)
	}
}

func vamp(attacker, defender *Obj) {
	a, d := attacker.Sheet, defender.Sheet
	vhp, vmp := d.MaxHP()/4, d.MaxMP()/4
	a.Heal(vhp)
	a.HealMP(vmp)
}
