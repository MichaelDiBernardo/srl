package client

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

type Console struct {
}

// Looks exactly like termbox.Cell :(
type glyph struct {
	Ch rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

var featureGlyphs = map[game.FeatureType]glyph{
	"FeatWall":  glyph{Ch: '#', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	"FeatFloor": glyph{Ch: '.', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
}

func NewConsole() *Console {
	return &Console{}
}

func (*Console) Init() error {
	return termbox.Init()
}

func (*Console) Render(w *game.World) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, row := range w.Level.Map {
		for _, tile := range row {
			if tile.Actor != nil {
				termbox.SetCell(tile.Pos.X, tile.Pos.Y, '@', termbox.ColorWhite, termbox.ColorBlack)
			} else {
				gl := featureGlyphs[tile.Feature.Type]
				termbox.SetCell(tile.Pos.X, tile.Pos.Y, gl.Ch, gl.Fg, gl.Bg)
			}
		}
	}
	termbox.Flush()
}

func (*Console) NextCommand() game.Command {
	keymap := map[rune]game.Command{
		'h': game.CommandMoveW,
		'j': game.CommandMoveS,
		'k': game.CommandMoveN,
		'l': game.CommandMoveE,
		'q': game.CommandQuit,
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
