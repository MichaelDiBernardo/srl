package console

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// A console client.
type Console struct {
	display display
	screens map[game.Mode]screen
	curscr  screen
}

// Create a new console client.
func New() *Console {
	display := &tbdisplay{}
	screens := map[game.Mode]screen{
		game.ModeHud:       newHudScreen(display),
		game.ModeInventory: newInventoryScreen(display),
		game.ModePickup:    newPickupScreen(display),
        game.ModeEquip:     newEquipScreen(display),
	}
	console := &Console{
		display: display,
		screens: screens,
	}
	console.switchScreen(game.ModeHud)
	return console
}

// Init the console client.
func (c *Console) Init() error {
	return c.display.Init()
}

// Render the current screen on this console.
func (c *Console) Render(g *game.Game) {
	c.display.Clear(termbox.ColorDefault, termbox.ColorDefault)
	c.curscr.Render(g)
	c.display.Flush()
}

// Handle an event generated by the game after the last command.
func (c *Console) Handle(ev game.Event) {
	switch e := ev.(type) {
	case game.ModeEvent:
		c.switchScreen(e.Mode)
	default:
		c.curscr.Handle(e)
	}
}

// Get the next command from the player to be sent to the game instance.
func (c *Console) NextCommand() game.Command {
	return c.curscr.NextCommand()
}

// Tear down the console client.
func (c *Console) Close() {
	c.display.Close()
}

// Switches screens in client to accomodate a change in game mode.
func (c *Console) switchScreen(mode game.Mode) {
	c.curscr = c.screens[mode]
}
