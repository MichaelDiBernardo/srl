package game

var (
	FeatWall       = &Feature{Type: "FeatWall", Solid: true, Opaque: true}
	FeatFloor      = &Feature{Type: "FeatFloor", Solid: false, Opaque: false}
	FeatClosedDoor = &Feature{Type: "FeatClosedDoor", Solid: true, Opaque: true}
	FeatOpenDoor   = &Feature{Type: "FeatOpenDoor", Solid: false, Opaque: true}
)
