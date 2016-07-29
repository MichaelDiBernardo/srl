package console

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/MichaelDiBernardo/srl/lib/math"
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

// Panel that renders the gameplay map.
type mapPanel struct {
	display display
}

// Where the map panel should go.
var mapPanelAnchor = math.Origin

// Create a new mapPanel.
func newMapPanel(display display) *mapPanel {
	return &mapPanel{display: display}
}

// Listens to nothing.
func (m *mapPanel) Handle(e game.Event) {
}

// Render the gameplay map to the hud.
func (m *mapPanel) Render(g *game.Game) {
	for _, row := range g.Level.Map {
		for _, tile := range row {
			pos := mapPanelAnchor.Add(tile.Pos)
			if tile.Actor != nil {
				gl := actorGlyphs[tile.Actor.Spec.Subtype]
				m.display.SetCell(pos.X, pos.Y, gl.Ch, gl.Fg, gl.Bg)
			} else {
				gl := featureGlyphs[tile.Feature.Type]
				m.display.SetCell(pos.X, pos.Y, gl.Ch, gl.Fg, gl.Bg)
			}
		}
	}
}

// The message panel, where messages are rendered at the bottom of the hud.
type messagePanel struct {
	display display
	lines   *list.List
	size    int
}

// Where the message panel should go.
var messagePanelAnchor = math.Pt(1, 17)

// Create a new messagePanel.
func newMessagePanel(size int, display display) *messagePanel {
	return &messagePanel{
		display: display,
		lines:   list.New(),
		size:    size,
	}
}

// Listens for new message events to build the message list.
func (m *messagePanel) Handle(e game.Event) {
	switch ev := e.(type) {
	case *game.MessageEvent:
		m.message(ev.Text)
	}
}

// Render the panel to the display.
func (m *messagePanel) Render(_ *game.Game) {
	for e, i := m.lines.Front(), 0; e != nil; e, i = e.Next(), i+1 {
		line := e.Value.(string)
		m.display.Write(messagePanelAnchor.X, messagePanelAnchor.Y+i, line, termbox.ColorWhite, termbox.ColorBlack)
	}
}

// Add a message to the list of messages to render.
func (m *messagePanel) message(text string) {
	m.lines.PushBack(text)
	if m.lines.Len() > m.size {
		m.lines.Remove(m.lines.Front())
	}
}
