package game

import (
	"github.com/MichaelDiBernardo/srl/lib/client"
	"github.com/MichaelDiBernardo/srl/lib/event"
	"github.com/MichaelDiBernardo/srl/lib/world"
)

type Game struct {
	client client.Client
	world  *world.World
}

func New() *Game {
	return &Game{
		client: client.NewConsole(),
		world:  world.New(),
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
		if ev == event.Quit {
			return
		}
		g.world.Handle(ev)
	}
}
