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

var actorGlyphs = map[game.ObjSubtype]glyph{
	"Player": glyph{Ch: '@', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
	"MonOrc": glyph{Ch: 'o', Fg: termbox.ColorGreen, Bg: termbox.ColorBlack},
}

var featureGlyphs = map[game.FeatureType]glyph{
	"FeatWall":  glyph{Ch: '#', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	"FeatFloor": glyph{Ch: '.', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
}

func NewConsole() *Console {
	return &Console{}
}

func (c *Console) Init() error {
	return termbox.Init()
}

func (c *Console) Render(g *game.Game) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, row := range g.Level.Map {
		for _, tile := range row {
			if tile.Actor != nil {
				gl := actorGlyphs[tile.Actor.Spec.Subtype]
				termbox.SetCell(tile.Pos.X, tile.Pos.Y, gl.Ch, gl.Fg, gl.Bg)
			} else {
				gl := featureGlyphs[tile.Feature.Type]
				termbox.SetCell(tile.Pos.X, tile.Pos.Y, gl.Ch, gl.Fg, gl.Bg)
			}
		}
	}

	line := 0
	for !g.Events.Empty() {
		m := g.Events.Next().(*game.MessageEvent)
		c.write(0, line, m.Text, termbox.ColorWhite, termbox.ColorBlack)
		line++
	}
	termbox.Flush()
}

func (c *Console) NextCommand() game.Command {
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

func (c *Console) Close() {
	termbox.Close()
}

func (c *Console) write(x, y int, text string, fg termbox.Attribute, bg termbox.Attribute) {
	i := 0
	for _, r := range text {
		termbox.SetCell(x+i, y, r, fg, bg)
		i++
	}
}
