package console

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

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

type messagePanel struct {
	lines *list.List
	size  int
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
