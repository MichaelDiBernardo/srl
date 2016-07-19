package client

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
)

type Renderer interface {
	Render(w *game.World)
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
