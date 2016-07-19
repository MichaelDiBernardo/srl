package game

type Event int

const (
	_ Event = iota
	EQuit
	EMoveN
	EMoveNE
    EMoveE
    EMoveSE
	EMoveS
	EMoveSW
	EMoveW
	EMoveNW
)
