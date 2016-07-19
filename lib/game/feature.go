package game

type FeatureType string

type Feature struct {
	Type   FeatureType
	Solid  bool
	Opaque bool
}

var FeatWall = &Feature{Type: "FeatWall", Solid: true, Opaque: true}
var FeatFloor = &Feature{Type: "FeatFloor", Solid: false, Opaque: true}
