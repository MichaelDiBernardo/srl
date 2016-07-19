package client

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
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

func (*Console) Render(w *game.World) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(w.Player.Pos.X, w.Player.Pos.Y, '@', termbox.ColorWhite, termbox.ColorBlack)
	termbox.Flush()
}

func (*Console) NextEvent() game.Event {
	keymap := map[rune]game.Event{
		'h': game.EMoveW,
		'j': game.EMoveS,
		'k': game.EMoveN,
		'l': game.EMoveE,
		'q': game.EQuit,
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
