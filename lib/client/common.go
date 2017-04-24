package client

import (
	"errors"
	"github.com/MichaelDiBernardo/srl/lib/game"
)

type Renderer interface {
	// Render this thing, using the given game as a read-only context. Do not
	// mutate this game instance!
	Render(g *game.Game)
}

type Initer interface {
	Init() error
}

type Closer interface {
	Close()
}

type Poller interface {
	// Wait for and handle user input, using the given game as context. Return
	// a command that should be issued to the game if one arises, or an error
	// ("ErrNoCommand") if we want to re-poll without pushing to the game.
	// The given game is intended to be _read-only_, do not mutate it.
	// Mutations should be done only through commands.
	Poll(g *game.Game) (game.Command, error)
}

var ErrNoCommand = errors.New("NoCommand")

type EventHandler interface {
	HandleEvent(game.Event)
}

type Client interface {
	Closer
	Initer
	Renderer
	Poller
	EventHandler
}
