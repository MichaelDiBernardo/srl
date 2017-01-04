package game

var (
	FeatWall       = &Feature{Type: "FeatWall", Solid: true, Opaque: true}
	FeatFloor      = &Feature{Type: "FeatFloor", Solid: false, Opaque: false}
	FeatClosedDoor = &Feature{Type: "FeatClosedDoor", Solid: false, Opaque: true}
)
