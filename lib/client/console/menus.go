package console

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// The inventory screen.
type inventoryScreen struct {
	display display
	menu    panel
}

// Create a new inventory screen.
func newInventoryScreen(display display) *inventoryScreen {
	return &inventoryScreen{
		display: display,
		menu:    newInventoryPanel(display),
	}
}

// Render the inventory screen.
func (inv *inventoryScreen) Render(g *game.Game) {
	inv.menu.Render(g)
}

// Handle an event generated by the game after the last command.
func (inv *inventoryScreen) Handle(ev game.Event) {
}

// Get the next command from the player to be sent to the game instance.
func (inv *inventoryScreen) NextCommand() game.Command {
	for {
		tboxev := inv.display.PollEvent()
		if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
			return game.ModeCommand{Mode: game.ModeHud}
		}
	}
}

// Panel that renders the inventory.
type inventoryPanel struct {
	display display
}

// Create a new inventoryPanel.
func newInventoryPanel(display display) *inventoryPanel {
	return &inventoryPanel{display: display}
}

// Listens to nothing.
func (m *inventoryPanel) Handle(e game.Event) {
}

// Render the menu.
func (m *inventoryPanel) Render(g *game.Game) {
	inv := g.Player.Packer.Inventory()
	renderInventory(m.display, "Inventory", inv)
}

// The pickup screen.
type pickupScreen struct {
	display display
	menu    panel
}

// Create a new pickup screen.
func newPickupScreen(display display) *pickupScreen {
	return &pickupScreen{
		display: display,
		menu:    newPickupPanel(display),
	}
}

// Render the pickup screen.
func (p *pickupScreen) Render(g *game.Game) {
	p.menu.Render(g)
}

// Handle an event generated by the game after the last command.
func (p *pickupScreen) Handle(ev game.Event) {
}

// Get the next command from the player to be sent to the game instance.
func (inv *pickupScreen) NextCommand() game.Command {
	for {
		tboxev := inv.display.PollEvent()
		if tboxev.Type != termbox.EventKey {
			continue
		}
		if tboxev.Key == termbox.KeyEsc {
			return game.ModeCommand{Mode: game.ModeHud}
		} else if ch := tboxev.Ch; ch != 0 {
			opt := selectOption(ch)
			if opt != -1 {
				return game.MenuCommand{Option: opt}
			}
		}
	}
}

// Panel that renders the pickup screen.
type pickupPanel struct {
	display display
}

// Create a new pickupPanel.
func newPickupPanel(display display) *pickupPanel {
	return &pickupPanel{display: display}
}

// Listens to nothing.
func (m *pickupPanel) Handle(e game.Event) {
}

// Render the menu.
func (m *pickupPanel) Render(g *game.Game) {
	inv := g.Player.Tile.Items
	renderInventory(m.display, "Ground", inv)
}

// The equip screen.
type equipScreen struct {
	display display
	menu    panel
}

// Create a new equip screen.
func newEquipScreen(display display) *equipScreen {
	return &equipScreen{
		display: display,
		menu:    newEquipPanel(display),
	}
}

// Render the inventory screen.
func (inv *equipScreen) Render(g *game.Game) {
	inv.menu.Render(g)
}

func (inv *equipScreen) Handle(ev game.Event) {
}

func (inv *equipScreen) NextCommand() game.Command {
	for {
		tboxev := inv.display.PollEvent()
		if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
			return game.ModeCommand{Mode: game.ModeHud}
		} else if ch := tboxev.Ch; ch != 0 {
			opt := selectOption(ch)
			if opt != -1 {
				return game.MenuCommand{Option: opt}
			}
		}
	}
}

// Panel that renders the equipment list.
type equipPanel struct {
	display display
}

// Create a new equipPanel.
func newEquipPanel(display display) *equipPanel {
	return &equipPanel{display: display}
}

// Listens to nothing.
func (m *equipPanel) Handle(e game.Event) {
}

// Render the menu.
func (m *equipPanel) Render(g *game.Game) {
	inv := g.Player.Packer.Inventory()
	renderInventory(m.display, "Equipment", inv)
}

// Function that renders a single inventory within an inventory menu panel.
func renderInventory(display display, title string, inv *game.Inventory) {
	items := inv.Items
	display.Write(0, 0, title, termbox.ColorWhite, termbox.ColorBlack)
	i := 0
	for e := items.Back(); e != nil; e = e.Prev() {
		item := e.Value.(*game.Obj)
		display.Write(1, 1+i, fmt.Sprintf("%c - %v", alphabet[i], item.Spec.Name), termbox.ColorWhite, termbox.ColorBlack)
		i++
	}
}

// Basic choosy things.
var alphabet = []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}

// Converts a selected menu option into an index that the game will use.
// Returns -1 if not found.
func selectOption(ch rune) int {
	for i, r := range alphabet {
		if r == ch {
			return i
		}
	}
	return -1
}
