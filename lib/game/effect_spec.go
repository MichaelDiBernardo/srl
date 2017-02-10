package game

const (
	BrandFire Effect = iota
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

var Brands = Effects{
	BrandFire,
	BrandElec,
	BrandIce,
	BrandPoison,
}

var ResistMap = map[Effect]Effect{
	BrandFire:    ResistFire,
	BrandElec:    ResistElec,
	BrandIce:     ResistIce,
	BrandPoison:  ResistPoison,
	EffectStun:   ResistStun,
	EffectPoison: ResistPoison,
}

var VulnMap = map[Effect]Effect{
	BrandFire:    VulnFire,
	BrandElec:    VulnElec,
	BrandIce:     VulnIce,
	BrandPoison:  VulnPoison,
	EffectStun:   VulnStun,
	EffectPoison: VulnPoison,
}

var EffectVerbs = map[Effect]string{
	BrandFire:   "burns",
	BrandElec:   "shocks",
	BrandIce:    "freezes",
	BrandPoison: "poisons",
}
