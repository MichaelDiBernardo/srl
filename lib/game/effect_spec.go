package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Brands, effects, resists, vulnerabilities found in the game.
const (
	// Sentinel.
	EffectNone Effect = iota

	// Brands
	BrandFire
	BrandElec
	BrandIce
	BrandPoison

	// Effects
	EffectBaseRegen // Regen that is applied every tick to every actor.
	EffectStun
	EffectPoison

	// Resists
	ResistFire
	ResistElec
	ResistIce
	ResistStun
	ResistPoison

	// Sentinel
	NumEffects
)

// Types of permissible effects.
const (
	EffectTypeBrand = iota
	EffectTypeResist
	EffectTypeStatus
)

// Mapping of effects to their types.
var EffectsSpecs = EffectsSpec{
	BrandFire:   {Type: EffectTypeBrand, ResistedBy: ResistFire, Verb: "burns"},
	BrandElec:   {Type: EffectTypeBrand, ResistedBy: ResistElec, Verb: "shocks"},
	BrandIce:    {Type: EffectTypeBrand, ResistedBy: ResistIce, Verb: "freezes"},
	BrandPoison: {Type: EffectTypeBrand, ResistedBy: ResistPoison, Verb: "poisons"},

	EffectBaseRegen: {Type: EffectTypeStatus},
	EffectStun:      {Type: EffectTypeStatus, ResistedBy: ResistStun},
	EffectPoison:    {Type: EffectTypeStatus, ResistedBy: ResistPoison},

	ResistFire:   {Type: EffectTypeResist},
	ResistElec:   {Type: EffectTypeResist},
	ResistIce:    {Type: EffectTypeResist},
	ResistStun:   {Type: EffectTypeResist},
	ResistPoison: {Type: EffectTypeResist},
}

// Prototype map for effects that are applied every tick.
var ActiveEffects = map[Effect]ActiveEffect{
	EffectBaseRegen: ActiveEffect{
		OnTick: regen,
	},
	EffectPoison: ActiveEffect{
		OnBegin: func(_ *ActiveEffect, t *ActorTicker) {
			t.Obj().Game.Events.Message(fmt.Sprintf("%s is poisoned.", t.Obj().Spec.Name))
		},
		OnTick: poison,
		OnEnd: func(_ *ActiveEffect, t *ActorTicker) {
			t.Obj().Game.Events.Message(fmt.Sprintf("%s recovers from poison.", t.Obj().Spec.Name))
		},
	},
}

// Regenerate the actor every turn.
func regen(e *ActiveEffect, t *ActorTicker, diff int) bool {
	sheet := t.obj.Sheet
	regen := sheet.Regen()

	e.Counter += regen * diff
	delayPerHp := RegenPeriod * GetDelay(2) / sheet.MaxHP()
	heal := e.Counter / delayPerHp

	if heal > 0 {
		sheet.Heal(heal)
		e.Counter -= heal * delayPerHp
	}
	return false
}

// We expect a speed 2 actor to fully recover in 100 turns.
const RegenPeriod = 100

// Apply poison damage each turn.
func poison(e *ActiveEffect, t *ActorTicker, _ int) bool {
	dmg := math.Max(20*e.Counter/100, 1)
	t.Obj().Sheet.Hurt(dmg)
	e.Counter -= dmg
	return e.Counter <= 0
}
