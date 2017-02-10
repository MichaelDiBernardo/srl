package game

// Brands, effects, resists, vulnerabilities found in the game.
const (
	// Regen that is applied every tick to every actor.
	EffectBaseRegen Effect = iota

	BrandFire
	BrandElec
	BrandIce
	BrandPoison
	EffectStun
	EffectPoison

	ResistFire
	ResistElec
	ResistIce
	ResistStun
	ResistPoison

	VulnFire
	VulnElec
	VulnIce
	VulnStun
	VulnPoison

	NumEffects
)

// The subset of effects that are brands.
var Brands = Effects{
	BrandFire,
	BrandElec,
	BrandIce,
	BrandPoison,
}

// Which effects are resisted by which.
var ResistMap = map[Effect]Effect{
	BrandFire:    ResistFire,
	BrandElec:    ResistElec,
	BrandIce:     ResistIce,
	BrandPoison:  ResistPoison,
	EffectStun:   ResistStun,
	EffectPoison: ResistPoison,
}

// Which effects are vulnerable to which.
var VulnMap = map[Effect]Effect{
	BrandFire:    VulnFire,
	BrandElec:    VulnElec,
	BrandIce:     VulnIce,
	BrandPoison:  VulnPoison,
	EffectStun:   VulnStun,
	EffectPoison: VulnPoison,
}

// Verbs associated with effects when they take effect in melee.
var EffectVerbs = map[Effect]string{
	BrandFire:   "burns",
	BrandElec:   "shocks",
	BrandIce:    "freezes",
	BrandPoison: "poisons",
}

// Prototype map for effects that are applied every tick.
var ActiveEffects = map[Effect]ActiveEffect{
	EffectBaseRegen: ActiveEffect{
		OnTick: regen,
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
