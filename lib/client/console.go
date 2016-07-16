package client

import (
	"github.com/MichaelDiBernardo/srl/lib/world"
	"github.com/nsf/termbox-go"
)

type Console struct {
}

func NewConsole() *Console {
	return &Console{}
}

func (*Console) Init() error {
	return termbox.Init()
}

func (*Console) Render(w *world.World) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(w.Player.X, w.Player.Y, '@', termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()
}

func (*Console) Close() {
	termbox.Close()
}
