package client

import (
	"github.com/MichaelDiBernardo/srl/lib/event"
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

func (*Console) NextEvent() event.Event {
	keymap := map[rune]event.Event{
		'h': event.MoveW,
		'j': event.MoveS,
		'k': event.MoveN,
		'l': event.MoveE,
		'q': event.Quit,
	}
	for {
		tboxev := termbox.PollEvent()

		if tboxev.Type != termbox.EventKey || tboxev.Key != 0 {
			continue
		}

		srlev := keymap[tboxev.Ch]
		if srlev != 0 {
			return srlev
		}
	}
}

func (*Console) Close() {
	termbox.Close()
}
