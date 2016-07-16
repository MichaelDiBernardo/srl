package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
    "github.com/MichaelDiBernardo/srl/lib/client"
    "github.com/MichaelDiBernardo/srl/lib/world"
    "github.com/nsf/termbox-go"
)

type Game struct {
    client client.Client
    world *world.World
}

func New() *Game {
    return &Game{
        client: client.NewConsole(),
        world: world.New(),
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

        player := g.world.Player
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowUp:
				player.Y = math.Max(player.Y-1, 0)
			case termbox.KeyArrowDown:
				player.Y += 1
			case termbox.KeyArrowLeft:
				player.X = math.Max(player.X-1, 0)
			case termbox.KeyArrowRight:
				player.X += 1
			}
		}
	}
}
