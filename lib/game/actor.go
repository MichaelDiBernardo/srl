package game

import (
    "github.com/MichaelDiBernardo/srl/lib/math"
)

type Player struct {
    Pos math.Point
}

func NewPlayer() *Player {
    return &Player{Pos: math.Origin}
}
