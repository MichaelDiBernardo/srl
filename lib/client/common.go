package client

import (
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

type Poller interface {
}

type Client interface {
    Renderer
    Initer
    Closer
}
