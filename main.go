package main

import (
	"github.com/MichaelDiBernardo/srl/lib/client"
	"github.com/MichaelDiBernardo/srl/lib/game"
)

type Game struct {
	client client.Client
	world  *game.World
}

func NewGame() *Game {
	return &Game{
		client: client.NewConsole(),
		world:  game.NewWorld(),
	}
}

func (g *Game) Loop() {
	err := g.client.Init()
	if err != nil {
		panic(err)
	}
	defer g.client.Close()

	for {
		g.client.Render(g.world)
		ev := g.client.NextEvent()
		if ev == game.EQuit {
			return
		}
		g.world.Handle(ev)
	}
}

func main() {
	game := NewGame()
	game.Loop()
}
