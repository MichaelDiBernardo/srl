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

type Eventer interface {
	NextEvent() game.Event
}

type Client interface {
	Closer
	Eventer
	Initer
	Renderer
}
