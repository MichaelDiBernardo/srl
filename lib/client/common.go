package client

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
)

type Renderer interface {
	Render(g *game.Game)
}

type Initer interface {
	Init() error
}

type Closer interface {
	Close()
}

type Poller interface {
	Poll() (game.Command, error)
}

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
