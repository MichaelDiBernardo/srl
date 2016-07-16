package actor

type Player struct {
    X, Y int
}

func NewPlayer() *Player {
    return &Player{X: 0, Y: 0}
}
