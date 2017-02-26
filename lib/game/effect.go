package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// An effect is something that a monster or a piece of equipment can have. This
// includes, brands, resists, status effects, etc.
type Effect uint

// An effect 'type' is a grouping of that effect e.g. is it a brand, a status
// effect, etc.
type EffectType uint

// Defines basic information about an effect -- i.e what type it is, and if it
// is an offensive effect or brand, what it is resisted by.
type EffectSpec struct {
	Type       EffectType
	ResistedBy Effect
	Verb       string
}

// Definition of _all_ effects and their types.
type EffectsSpec map[Effect]*EffectSpec

// A single instance of an effect on a creature, attack, defense, etc.
type EffectInfo struct {
	*EffectSpec
	Count int
}

// Given an effect and its count, this will create an EffectInfo that
// represents the same.
func newEffectInfo(effect Effect, count int) EffectInfo {
	spec := EffectsSpecs[effect]
	if spec == nil {
		panic(fmt.Sprintf("Could not find spec for effect %v", effect))
	}
	return EffectInfo{EffectSpec: spec, Count: count}
}

// A collection of effects on a monster, piece of equipment, etc.
type Effects map[Effect]EffectInfo

// Given a histogram of Effect -> resist/vuln count, this will create a Effects
// collection with the given counts in place.
func NewEffects(counts map[Effect]int) Effects {
	effects := Effects{}

	for effect, count := range counts {
		effects[effect] = newEffectInfo(effect, count)
	}

	return effects
}

// Returns the number of 'pips' of this effect that this collection has.
func (effects Effects) Has(effect Effect) int {
	return effects[effect].Count
}

// How many pips of resistance do I have against 'effect'? This can be
// negative, which indicates a vulnerability to 'effect'.
func (effects Effects) Resists(effect Effect) int {
	info := EffectsSpecs[effect]
	if info == nil {
		return 0
	}
	return effects.Has(info.ResistedBy)
}

// Filters out the brands from this collection of effects.
func (effects Effects) Brands() Effects {
	brands := make(Effects)
	for effect, info := range effects {
		if info.Type == EffectTypeBrand {
			brands[effect] = info
		}
	}
	return brands
}

// Produces a new Effects that is the union of 'e1' and 'e2', accumulating the
// counts of each effect. This does not mutate either of the inputs.
func (e1 Effects) Merge(e2 Effects) Effects {
	merged := Effects{}
	for k, v := range e1 {
		merged[k] = v
	}

	for k, v := range e2 {
		info := merged[k]
		v.Count += info.Count
		merged[k] = v
	}
	return merged
}

// Given the amount of raw damage done by an effect 'effect', this figures out
// how much damage should actually be done after resistances or vulnerabilities
// to 'effect' are taken into account.
func (e Effects) ResistDmg(effect Effect, dmg int) int {
	resists := e.Resists(effect)
	if resists >= 0 {
		return dmg / (resists + 1)
	} else {
		return dmg * (-resists + 1)
	}
}

// An instance of an effect that is currently affected an actor. These are
// managed by the actor's ticker.
type ActiveEffect struct {
	// The effect counter. Might be turns, accumulated delay, residual damage, etc.
	Counter int
	// What to do when the effect is first inflicted on an actor. The last
	// integer argument will hold the previous value of the counter for this
	// effect if it is being continued / augmented (e.g. "X is more poisoned")
	OnBegin func(*ActiveEffect, Ticker, int)
	// Responsible for updating Left given the delay diff, plus enforcing the
	// effect. A return value of 'true' indicates that the effect should be
	// terminated.
	OnTick func(*ActiveEffect, Ticker, int) bool
	// What to do when the effect has run its course.
	OnEnd func(*ActiveEffect, Ticker)
}

// Creates a new ActiveEffect record for the given effect.
func NewActiveEffect(e Effect, counter int) *ActiveEffect {
	ae, ok := &ActiveEffect{}, true
	*ae, ok = ActiveEffects[e]

	if !ok {
		panic(fmt.Sprintf("%v not found in ActiveEffects", e))
	}

	ae.Counter = counter

	if ae.OnBegin == nil {
		ae.OnBegin = func(*ActiveEffect, Ticker, int) {}
	}
	if ae.OnEnd == nil {
		ae.OnEnd = func(*ActiveEffect, Ticker) {}
	}
	return ae
}

// Implementation of specific active effects.
// TODO: This could have more simply been an interface type with a factory
// function to create the right implementation for the given effect enum value.
// I'm not quite sure how we got here. Maybe fix this.  Oops. Although, all
// those structs would really have looked the same and just been aliases of a
// base struct.
var (
	// Base regen that actors get every turn.
	AEBaseRegen = ActiveEffect{
		OnTick: func(e *ActiveEffect, t Ticker, diff int) bool {
			sheet := t.Obj().Sheet
			regen := sheet.Regen()

			e.Counter += regen * diff
			delayPerHp := RegenPeriod * GetDelay(2) / sheet.MaxHP()
			heal := e.Counter / delayPerHp

			if heal > 0 {
				sheet.Heal(heal)
				e.Counter -= heal * delayPerHp
			}
			return false
		},
	}
	AEPoison = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, prev int) {
			var msg string
			if prev == 0 {
				msg = "%s is poisoned."
			} else {
				msg = "%s is more poisoned."
			}
			t.Obj().Game.Events.Message(fmt.Sprintf(msg, t.Obj().Spec.Name))
		},
		OnTick: hpdecay,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Game.Events.Message(fmt.Sprintf("%s recovers from poison.", t.Obj().Spec.Name))
		},
	}
	AECut = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, prev int) {
			var msg string
			if prev == 0 {
				msg = "%s is wounded."
			} else {
				msg = "%s is more wounded."
			}
			t.Obj().Game.Events.Message(fmt.Sprintf(msg, t.Obj().Spec.Name))
		},
		OnTick: hpdecay,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Game.Events.Message(fmt.Sprintf("%s is healed from wounds.", t.Obj().Spec.Name))
		},
	}
	AEStun = ActiveEffect{
		OnBegin: func(e *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetStun(getstunlevel(e.Counter))
		},
		OnTick: func(e *ActiveEffect, t Ticker, _ int) bool {
			e.Counter -= 1
			t.Obj().Sheet.SetStun(getstunlevel(e.Counter))
			return e.Counter <= 0
		},
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetStun(NotStunned)
		},
	}
	AEBlind = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetBlind(true)
			msg := fmt.Sprintf("%s is blinded.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetBlind(false)
			msg := fmt.Sprintf("%s can see again.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AESlow = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetSlow(true)
			msg := fmt.Sprintf("%s is slowed.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetSlow(false)
			msg := fmt.Sprintf("%s speeds up again.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AEConfuse = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetConfused(true)
			msg := fmt.Sprintf("%s is confused.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetConfused(false)
			msg := fmt.Sprintf("%s recovers from confusion.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AEFear = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetAfraid(true)
			msg := fmt.Sprintf("%s is scared.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetAfraid(false)
			msg := fmt.Sprintf("%s recovers from fear.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AEPara = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetParalyzed(true)
			msg := fmt.Sprintf("%s is paralyzed.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetParalyzed(false)
			msg := fmt.Sprintf("%s can move again.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AESilence = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetSilenced(true)
			msg := fmt.Sprintf("%s is silenced.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetSilenced(false)
			msg := fmt.Sprintf("%s can speak again.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AECurse = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, _ int) {
			t.Obj().Sheet.SetCursed(true)
			msg := fmt.Sprintf("%s is cursed.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			t.Obj().Sheet.SetCursed(false)
			msg := fmt.Sprintf("%s is no longer cursed.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
		},
	}
	AEStim = ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t Ticker, prev int) {
			msg := fmt.Sprintf("%s feels a rush!", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
			if prev == 0 {
				modAllSkills(t.Obj().Sheet, 2)
			}
		},
		OnTick: basictick,
		OnEnd: func(_ *ActiveEffect, t Ticker) {
			msg := fmt.Sprintf("%s feels the rush wear off.", t.Obj().Spec.Name)
			t.Obj().Game.Events.Message(msg)
			modAllSkills(t.Obj().Sheet, -2)
		},
	}
)

// An actor's stun level depends on how many turns of stun they've accumulated.
func getstunlevel(cstun int) StunLevel {
	switch {
	case cstun <= 0:
		return NotStunned
	case (0 < cstun) && (cstun < 50):
		return Stunned
	default:
		return MoreStunned
	}
}

// Hurt a poisoned / cut actor.
func hpdecay(e *ActiveEffect, t Ticker, _ int) bool {
	dmg := math.Max(20*e.Counter/100, 1)
	t.Obj().Sheet.Hurt(dmg)
	e.Counter -= dmg
	return e.Counter <= 0
}

func basictick(e *ActiveEffect, t Ticker, _ int) bool {
	e.Counter -= 1
	return e.Counter <= 0
}

// We expect a speed 2 actor to fully recover in 100 turns.
const RegenPeriod = 100
