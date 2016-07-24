package game

type FeatureType string

type Feature struct {
	Type   FeatureType
	Solid  bool
	Opaque bool
}
