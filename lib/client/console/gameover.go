package console

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
	"strings"
)

func newGameOverScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newGameOverPanel(display)},
	}
}

type gameOverPanel struct {
	display display
}

func newGameOverPanel(display display) *gameOverPanel {
	return &gameOverPanel{display: display}
}

func (g *gameOverPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && (tboxev.Key == termbox.KeyEsc || tboxev.Key == termbox.KeyEnter) {
		return game.QuitCommand{}, nil
	}
	return nocommand()
}

// Listens to nothing.
func (p *gameOverPanel) HandleEvent(e game.Event) {
}

// Render the panel.
func (p *gameOverPanel) Render(g *game.Game) {
	msg := center(fmt.Sprintf("✝✝✝ Ur dead %s ✝✝✝", g.Player.Spec.Name), consoleBounds.Width(), " ")
	p.display.Write(0, 10, msg, termbox.ColorRed, termbox.ColorBlack)
}

func center(s string, width int, fill string) string {
	padding := (width - len([]rune(s))) / 2
	return strings.Repeat(fill, padding) + s + strings.Repeat(fill, padding)
}
