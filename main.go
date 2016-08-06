package main

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/client"
	"github.com/MichaelDiBernardo/srl/lib/client/console"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"io"
	"log"
	"os"
)

// A single running game. Once we get to serverland, srl will handle multiple
// of these simultaneously.
type Session struct {
	client client.Client
	game   *game.Game
}

func NewSession() *Session {
	g := game.NewGame()
	g.Start()
	return &Session{
		client: console.New(),
		game:   g,
	}
}

func (s *Session) Loop() {
	err := s.client.Init()
	if err != nil {
		panic(err)
	}
	defer s.client.Close()

	for {
		s.client.Render(s.game)
		command := s.client.NextCommand()
		_, quit := command.(game.QuitCommand)
		if quit {
			return
		}
		s.game.Handle(command)

		for !s.game.Events.Empty() {
			ev := s.game.Events.Next()
			s.client.Handle(ev)
		}
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

	s := NewSession()
	s.Loop()
}
