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

type Commander interface {
	NextCommand() game.Command
}

type Client interface {
	Closer
	Initer
	Renderer
	Commander
}
