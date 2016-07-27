package console

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// A console client.
type Console struct {
	display display
}

// Create a new console client.
func New() *Console {
	return &Console{display: &tbdisplay{}}
}

// Init the console client.
func (c *Console) Init() error {
	return c.display.Init()
}

// Render the current screen on this console.
func (c *Console) Render(g *game.Game) {
	c.display.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for _, row := range g.Level.Map {
		for _, tile := range row {
			if tile.Actor != nil {
				gl := actorGlyphs[tile.Actor.Spec.Subtype]
				c.display.SetCell(tile.Pos.X, tile.Pos.Y, gl.Ch, gl.Fg, gl.Bg)
			} else {
				gl := featureGlyphs[tile.Feature.Type]
				c.display.SetCell(tile.Pos.X, tile.Pos.Y, gl.Ch, gl.Fg, gl.Bg)
			}
		}
	}

	line := 0
	for !g.Events.Empty() {
		m := g.Events.Next().(*game.MessageEvent)
		c.display.Write(0, line, m.Text, termbox.ColorWhite, termbox.ColorBlack)
		line++
	}
	c.display.Flush()
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
		tboxev := c.display.PollEvent()

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
	c.display.Close()
}
