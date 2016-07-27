package client

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// A console client.
type Console struct {
}

// A glyph used to render a tile.
type glyph struct {
	Ch rune
	Fg termbox.Attribute
	Bg termbox.Attribute
}

// Glyphs used to render actors.
var actorGlyphs = map[game.ObjSubtype]glyph{
	"Player": glyph{Ch: '@', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
	"MonOrc": glyph{Ch: 'o', Fg: termbox.ColorGreen, Bg: termbox.ColorBlack},
}

// Glyphs used to render tiles.
var featureGlyphs = map[game.FeatureType]glyph{
	"FeatWall":  glyph{Ch: '#', Fg: termbox.ColorRed, Bg: termbox.ColorBlack},
	"FeatFloor": glyph{Ch: '.', Fg: termbox.ColorWhite, Bg: termbox.ColorBlack},
}

// Create a new console client.
func NewConsole() *Console {
	return &Console{}
}

// Init the console client.
func (c *Console) Init() error {
	return termbox.Init()
}

// Render the current screen on this console.
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
		write(0, line, m.Text, termbox.ColorWhite, termbox.ColorBlack)
		line++
	}
	termbox.Flush()
}

// Get the next command from the player to be sent to the game instance.
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

// Tear down the console client.
func (c *Console) Close() {
	termbox.Close()
}

type messagePanel struct {
	lines *list.List
    size int
}

func newMessagePanel(size int) *messagePanel {
    return &messagePanel{lines: list.New(), size: size}
}

func (m *messagePanel) Message(text string) {
    m.lines.PushBack(text)
    if m.lines.Len() > m.size {
        m.lines.Remove(m.lines.Front())
    }
}

// Write a string to the console.
func write(x, y int, text string, fg termbox.Attribute, bg termbox.Attribute) {
	i := 0
	for _, r := range text {
		termbox.SetCell(x+i, y, r, fg, bg)
		i++
	}
}
