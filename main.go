package main

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/client"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"io"
	"log"
	"os"
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
		command := g.client.NextCommand()
		if command == game.CommandQuit {
			return
		}
		g.world.Handle(command)
	}
}

var logfile *os.File

func setup() {
	f, err := os.OpenFile("srl.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening logfile: %v", err))
	}
	log.SetOutput(io.Writer(f))
}

func teardown() {
	logfile.Close()
}

func main() {
	setup()
	defer teardown()

	game := NewGame()
	game.Loop()
}
