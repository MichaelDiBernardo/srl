package game

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
	EffectCut

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
	EffectCut:       {Type: EffectTypeStatus},

	ResistFire:   {Type: EffectTypeResist},
	ResistElec:   {Type: EffectTypeResist},
	ResistIce:    {Type: EffectTypeResist},
	ResistStun:   {Type: EffectTypeResist},
	ResistPoison: {Type: EffectTypeResist},
}

// Prototype map for effects that are applied every tick.
var ActiveEffects = map[Effect]ActiveEffect{
	EffectBaseRegen: AEBaseRegen,
	EffectPoison:    AEPoison,
	EffectStun:      AEStun,
	EffectCut:       AECut,
}
