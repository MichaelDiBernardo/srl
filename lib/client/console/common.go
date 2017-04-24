package console

import (
	"github.com/MichaelDiBernardo/srl/lib/client"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// A screen is composed of multiple panels. It's responsible for polling the
// player for input, figuring out if any of the subpanels have responded to
// that input in a way that merits notifying the game (by producing a command),
// and then doing so.
type screen struct {
	display display
	panels  []panel
}

func (s *screen) Render(g *game.Game) {
	for _, p := range s.panels {
		p.Render(g)
	}
}

// Polls the player for input, and then asks every panel to handle it. If a
// panel responds back with a command that should be sent to the game, the
// other panels are not asked.
func (s *screen) Poll(g *game.Game) (game.Command, error) {
	tboxev := s.display.PollEvent()
	for _, p := range s.panels {
		command, err := p.HandleInput(g, tboxev)
		if err == nil {
			return command, err
		}
	}
	return game.NoCommand{}, client.ErrNoCommand
}

// Asks the panels to handle ev, in order.
func (s *screen) HandleEvent(ev game.Event) {
	for _, p := range s.panels {
		p.HandleEvent(ev)
	}
}

// A panel is a specific region on a screen. Each can handle its own subset of
// inputs and events it is interested in. Panels do not operate in their own
// coordinate space (yet) -- you have to be careful where you draw with
// Render(), because there is nothing stopping you from drawing over another
// panel's real estate.
type panel interface {
	client.Renderer
	client.EventHandler
	InputHandler
}

type InputHandler interface {
	// Respond to a user action, using the given game as read-only context if
	// required. Do not mutate this instance of the game! Returns a command if
	// one should be issued to the game, or an error if this thing has no
	// command to furnish.
	HandleInput(*game.Game, termbox.Event) (game.Command, error)
}

// Shortcut for saying "I don't have a command in response to this input."
func nocommand() (game.Command, error) {
	return game.NoCommand{}, client.ErrNoCommand
}
