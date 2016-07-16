package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"github.com/nsf/termbox-go"
)

type Game struct {
}

func New() *Game {
    return &Game{}
}

func (g *Game) Loop() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	x, y := 0, 0

	for {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x, y, '@', termbox.ColorBlue, termbox.ColorBlack)
		termbox.Flush()

		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
			case termbox.KeyArrowUp:
				y = math.Max(y-1, 0)
			case termbox.KeyArrowDown:
				y += 1
			case termbox.KeyArrowLeft:
				x = math.Max(x-1, 0)
			case termbox.KeyArrowRight:
				x += 1
			}
		}
	}
}
