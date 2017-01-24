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
)

var Brands = Effects{
	BrandFire,
	BrandElec,
	BrandIce,
}

var ResistMap = map[Effect]Effect{
	BrandFire:    ResistFire,
	BrandElec:    ResistElec,
	BrandIce:     ResistIce,
	BrandPoison:  ResistPoison,
	EffectStun:   ResistStun,
	EffectPoison: ResistPoison,
}

var EffectVerbs = map[Effect]string{
	BrandFire:   "burns",
	BrandElec:   "shocks",
	BrandIce:    "freezes",
	BrandPoison: "poisons",
}
