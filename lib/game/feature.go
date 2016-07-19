package game

type FeatureType string

type Feature struct {
    Type FeatureType
    Solid bool
    Opaque bool
}

var FeatWall = &Feature{Type: "FtWall", Solid: true, Opaque: true}
var FeatFloor = &Feature{Type: "FtFloor", Solid: false, Opaque: true}
