package event

type Event int

const (
	_ Event = iota
	Quit
	MoveN
	MoveNE
	MoveE
	MoveSE
	MoveS
	MoveSW
	MoveW
	MoveNW
)
