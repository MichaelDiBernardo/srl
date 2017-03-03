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
	BrandAcid

	// Slays
	SlayPearl
	SlayHunter
	SlayBattle
	SlayDispel

	// Effects
	EffectBaseRegen // Regen that is applied every tick to every actor.
	EffectStun
	EffectPoison
	EffectCut
	EffectBlind
	EffectSlow
	EffectConfuse
	EffectFear
	EffectPara
	EffectSilence
	EffectCurse
	EffectPetrify
	EffectBless
	EffectStim
	EffectHyper
	EffectVamp
	EffectShatter
	EffectDrainStr
	EffectDrainAgi
	EffectDrainVit
	EffectDrainMnd

	// Resists
	ResistFire
	ResistElec
	ResistIce
	ResistStun
	ResistPoison
	ResistBlind
	ResistSlow
	ResistConfuse
	ResistFear
	ResistPara
	ResistSilence
	ResistCurse
	ResistPetrify
	ResistCrit
	ResistVamp
	ResistAcid
	ResistDrain

	// Flags
	WeakPearl
	WeakHunter
	WeakBattle
	WeakDispel

	// Sentinel
	NumEffects
)

// Types of permissible effects.
const (
	EffectTypeBrand = iota
	EffectTypeResist
	EffectTypeStatus
	EffectTypeSlay
	EffectTypeFlag
)

// Mapping of effects to their types.
var EffectsSpecs = EffectsSpec{
	BrandFire:   {Type: EffectTypeBrand, ResistedBy: ResistFire, Verb: "burns"},
	BrandElec:   {Type: EffectTypeBrand, ResistedBy: ResistElec, Verb: "shocks"},
	BrandIce:    {Type: EffectTypeBrand, ResistedBy: ResistIce, Verb: "freezes"},
	BrandPoison: {Type: EffectTypeBrand, ResistedBy: ResistPoison, Verb: "poisons"},
	BrandAcid:   {Type: EffectTypeBrand, ResistedBy: ResistAcid, Verb: "melts"},

	SlayPearl:  {Type: EffectTypeSlay, Slays: WeakPearl},
	SlayHunter: {Type: EffectTypeSlay, Slays: WeakHunter},
	SlayBattle: {Type: EffectTypeSlay, Slays: WeakBattle},
	SlayDispel: {Type: EffectTypeSlay, Slays: WeakDispel},

	EffectBaseRegen: {Type: EffectTypeStatus},
	EffectStun:      {Type: EffectTypeStatus, ResistedBy: ResistStun},
	EffectPoison:    {Type: EffectTypeStatus, ResistedBy: ResistPoison},
	EffectCut:       {Type: EffectTypeStatus},
	EffectBlind:     {Type: EffectTypeStatus, ResistedBy: ResistBlind},
	EffectSlow:      {Type: EffectTypeStatus, ResistedBy: ResistSlow},
	EffectConfuse:   {Type: EffectTypeStatus, ResistedBy: ResistConfuse},
	EffectFear:      {Type: EffectTypeStatus, ResistedBy: ResistFear},
	EffectPara:      {Type: EffectTypeStatus, ResistedBy: ResistPara},
	EffectSilence:   {Type: EffectTypeStatus, ResistedBy: ResistSilence},
	EffectCurse:     {Type: EffectTypeStatus, ResistedBy: ResistCurse},
	EffectPetrify:   {Type: EffectTypeStatus, ResistedBy: ResistPetrify},
	EffectBless:     {Type: EffectTypeStatus},
	EffectStim:      {Type: EffectTypeStatus},
	EffectHyper:     {Type: EffectTypeStatus},
	EffectVamp:      {Type: EffectTypeStatus, ResistedBy: ResistVamp},
	EffectShatter:   {Type: EffectTypeStatus},
	EffectDrainStr:  {Type: EffectTypeStatus},
	EffectDrainAgi:  {Type: EffectTypeStatus},
	EffectDrainVit:  {Type: EffectTypeStatus},
	EffectDrainMnd:  {Type: EffectTypeStatus},

	WeakPearl:  {Type: EffectTypeFlag},
	WeakHunter: {Type: EffectTypeFlag},
	WeakBattle: {Type: EffectTypeFlag},
	WeakDispel: {Type: EffectTypeFlag},

	ResistFire:    {Type: EffectTypeResist},
	ResistElec:    {Type: EffectTypeResist},
	ResistIce:     {Type: EffectTypeResist},
	ResistAcid:    {Type: EffectTypeResist},
	ResistStun:    {Type: EffectTypeResist},
	ResistPoison:  {Type: EffectTypeResist},
	ResistBlind:   {Type: EffectTypeResist},
	ResistSlow:    {Type: EffectTypeResist},
	ResistConfuse: {Type: EffectTypeResist},
	ResistFear:    {Type: EffectTypeResist},
	ResistPara:    {Type: EffectTypeResist},
	ResistCurse:   {Type: EffectTypeResist},
	ResistPetrify: {Type: EffectTypeResist},
	ResistCrit:    {Type: EffectTypeResist},
	ResistVamp:    {Type: EffectTypeResist},
	ResistDrain:   {Type: EffectTypeResist},
}

// Prototype map for effects that are applied every tick.
var ActiveEffects = map[Effect]ActiveEffect{
	EffectBaseRegen: AEBaseRegen,
	EffectPoison:    AEPoison,
	EffectStun:      AEStun,
	EffectCut:       AECut,
	EffectBlind:     AEBlind,
	EffectSlow:      AESlow,
	EffectConfuse:   AEConfuse,
	EffectFear:      AEFear,
	EffectPara:      AEPara,
	EffectSilence:   AESilence,
	EffectCurse:     AECurse,
	EffectStim:      AEStim,
	EffectHyper:     AEHyper,
	EffectPetrify:   AEPetrify,
	EffectShatter:   AECorrode,
	EffectDrainStr:  AEDrainStr,
	EffectDrainAgi:  AEDrainAgi,
	EffectDrainVit:  AEDrainVit,
	EffectDrainMnd:  AEDrainMnd,
}
