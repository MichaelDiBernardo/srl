package console

import (
	"container/list"
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"github.com/nsf/termbox-go"
)

// The HUD screen.
type hudScreen struct {
	display display
	panels  [3]panel
}

var hudBounds = math.Rect(math.Origin, math.Pt(80, 24))

// Where the message panel should go.
var messagePanelBounds = math.Rect(math.Pt(1, 19), math.Pt(1, hudBounds.Max.Y))

// How many lines should show by default in the message panel.
var messagePanelNumLines = messagePanelBounds.Height()

// Where the status panel should go.
var statusPanelBounds = math.Rect(math.Pt(38, 0), math.Pt(hudBounds.Max.X-38, messagePanelBounds.Min.Y))

// Where the map panel should go.
var mapPanelBounds = math.Rect(math.Origin, math.Pt(statusPanelBounds.Min.X, messagePanelBounds.Min.Y))

// Create a new HUD.
func newHudScreen(display display) *hudScreen {
	return &hudScreen{
		display: display,
		panels: [3]panel{
			newMapPanel(display),
			newMessagePanel(messagePanelNumLines, display),
			newStatusPanel(display),
		},
	}
}

// Render the HUD.
func (h *hudScreen) Render(g *game.Game) {
	for _, p := range h.panels {
		p.Render(g)
	}
}

// Handle an event generated by the game after the last command.
func (h *hudScreen) Handle(ev game.Event) {
	for _, p := range h.panels {
		p.Handle(ev)
	}
}

// Get the next command from the player to be sent to the game instance.
func (h *hudScreen) NextCommand() game.Command {
	keymap := map[rune]game.Command{
		'h': game.CommandMoveW,
		'j': game.CommandMoveS,
		'k': game.CommandMoveN,
		'l': game.CommandMoveE,
		'q': game.CommandQuit,
	}
	for {
		tboxev := h.display.PollEvent()

		if tboxev.Type != termbox.EventKey || tboxev.Key != 0 {
			continue
		}

		srlev := keymap[tboxev.Ch]
		if srlev != 0 {
			return srlev
		}
	}
}

// Panel that renders the gameplay map.
type mapPanel struct {
	display display
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

// Create a new mapPanel.
func newMapPanel(display display) *mapPanel {
	return &mapPanel{display: display}
}

// Listens to nothing.
func (m *mapPanel) Handle(e game.Event) {
}

// Render the gameplay map to the hud.
func (m *mapPanel) Render(g *game.Game) {
	center := g.Player.Pos()
	boundsdist := math.Pt(mapPanelBounds.Width()/2, mapPanelBounds.Height()/2)
	viewport := math.Rect(center.Sub(boundsdist), center.Add(boundsdist))
	maptrans := mapPanelBounds.Min.Sub(viewport.Min)
	level := g.Level

	for x := viewport.Min.X; x < viewport.Max.X; x++ {
		for y := viewport.Min.Y; y < viewport.Max.Y; y++ {
			cur := math.Pt(x, y)
			if !cur.In(level) {
				continue
			}

			tile := level.At(cur)
			drawpos := cur.Add(maptrans)

			if tile.Actor != nil {
				gl := actorGlyphs[tile.Actor.Spec.Subtype]
				m.display.SetCell(drawpos.X, drawpos.Y, gl.Ch, gl.Fg, gl.Bg)
			} else {
				gl := featureGlyphs[tile.Feature.Type]
				m.display.SetCell(drawpos.X, drawpos.Y, gl.Ch, gl.Fg, gl.Bg)
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
		m.display.Write(messagePanelBounds.Min.X, messagePanelBounds.Min.Y+i, line, termbox.ColorWhite, termbox.ColorBlack)
	}
}

// Add a message to the list of messages to render.
func (m *messagePanel) message(text string) {
	m.lines.PushBack(text)
	if m.lines.Len() > m.size {
		m.lines.Remove(m.lines.Front())
	}
}

// Panel that renders player status on hud.
type statusPanel struct {
	display display
}

// Create a new statsPanel.
func newStatusPanel(display display) *statusPanel {
	return &statusPanel{display: display}
}

// Listens for nothing.
func (s *statusPanel) Handle(e game.Event) {
}

// Render the panel to the display.
func (s *statusPanel) Render(g *game.Game) {
	player := g.Player
	fg, bg := termbox.ColorWhite, termbox.ColorBlack
	stats := player.Stats
	sheet := player.Sheet

	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+0, player.Spec.Name, fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+1, "Human", fg, bg)

	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+3, fmt.Sprintf("%-7s%3d", "STR", stats.Str()), fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+4, fmt.Sprintf("%-7s%3d", "AGI", stats.Agi()), fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+5, fmt.Sprintf("%-7s%3d", "VIT", stats.Vit()), fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+6, fmt.Sprintf("%-7s%3d", "MND", stats.Mnd()), fg, bg)

	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+8, fmt.Sprintf("%-7s%3d:%-3d", "HP", sheet.HP(), sheet.MaxHP()), fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+9, fmt.Sprintf("%-7s%3d:%-3d", "MP", sheet.MP(), sheet.MaxMP()), fg, bg)

	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+11, fmt.Sprintf("%-7s%8s", "FIGHT", "(+2,2d6)"), fg, bg)
	s.display.Write(statusPanelBounds.Min.X, statusPanelBounds.Min.Y+12, fmt.Sprintf("%-7s%8s", "DEF", "[+0,1-3]"), fg, bg)
}
