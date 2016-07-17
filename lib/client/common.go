package client

import (
	"github.com/MichaelDiBernardo/srl/lib/event"
	"github.com/MichaelDiBernardo/srl/lib/world"
)

type Renderer interface {
	Render(w *world.World)
}

type Initer interface {
	Init() error
}

type Closer interface {
	Close()
}

type Eventer interface {
	NextEvent() event.Event
}

type Client interface {
	Closer
	Eventer
	Initer
	Renderer
}
