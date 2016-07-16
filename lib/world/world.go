package world

import "github.com/MichaelDiBernardo/srl/lib/actor"

type World struct {
    Player *actor.Player
}

func New() *World {
    return &World{Player: actor.NewPlayer()}
}
