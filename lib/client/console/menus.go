package console

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/game"
	"github.com/nsf/termbox-go"
)

// Create a new inventory screen.
func newInventoryScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newInventoryPanel(display)},
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

func (inv *inventoryPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	}
	return nocommand()
}

// Listens to nothing.
func (m *inventoryPanel) HandleEvent(e game.Event) {
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
func newPickupScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newPickupPanel(display)},
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

func (p *pickupPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type != termbox.EventKey {
		return nocommand()
	}
	if tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	} else if ch := tboxev.Ch; ch != 0 {
		opt := selectOption(ch)
		if opt != -1 {
			return game.MenuCommand{Option: opt}, nil
		}
	}
	return nocommand()
}

// Listens to nothing.
func (p *pickupPanel) HandleEvent(e game.Event) {
}

// Render the menu.
func (p *pickupPanel) Render(g *game.Game) {
	inv := g.Player.Tile.Items
	renderInventory(p.display, "Take what?", inv)
}

// Create a new drop screen.
func newDropScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newDropPanel(display)},
	}
}

// Panel that renders the drop screen.
type dropPanel struct {
	display display
}

// Create a new dropPanel.
func newDropPanel(display display) *dropPanel {
	return &dropPanel{display: display}
}

func (d *dropPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type != termbox.EventKey {
		return nocommand()
	}
	if tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	} else if ch := tboxev.Ch; ch != 0 {
		opt := selectOption(ch)
		if opt != -1 {
			return game.MenuCommand{Option: opt}, nil
		}
	}
	return nocommand()
}

// Listens to nothing.
func (d *dropPanel) HandleEvent(e game.Event) {
}

// Render the menu.
func (d *dropPanel) Render(g *game.Game) {
	inv := g.Player.Packer.Inventory()
	renderInventory(d.display, "Drop what?", inv)
}

// Create a new equip screen.
func newEquipScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newEquipPanel(display)},
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

func (e *equipPanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	} else if ch := tboxev.Ch; ch != 0 {
		opt := selectOption(ch)
		if opt != -1 {
			return game.MenuCommand{Option: opt}, nil
		}
	}
	return nocommand()
}

// Listens to nothing.
func (e *equipPanel) HandleEvent(ev game.Event) {
}

// Render the menu.
func (e *equipPanel) Render(g *game.Game) {
	inv := g.Player.Packer.Inventory()
	renderInventory(e.display, "Equip what?", inv)
}

// Create a new remove screen.
func newRemoveScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newRemovePanel(display)},
	}
}

// Panel that renders the remove list.
type removePanel struct {
	display display
}

// Create a new removePanel.
func newRemovePanel(display display) *removePanel {
	return &removePanel{display: display}
}

func (r *removePanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	} else if ch := tboxev.Ch; ch != 0 {
		opt := selectOption(ch)
		if opt != -1 {
			return game.MenuCommand{Option: opt}, nil
		}
	}
	return nocommand()
}

// Listens to nothing.
func (r *removePanel) HandleEvent(e game.Event) {
}

// Render the menu.
func (r *removePanel) Render(g *game.Game) {
	display := r.display
	body := g.Player.Equipper.Body()

	display.Write(0, 0, "Remove what?", termbox.ColorWhite, termbox.ColorBlack)

	for i, equip := range body.Slots {
		name := "(nothing)"
		if equip != nil {
			name = equip.Spec.Name
		}
		display.Write(1, 1+i, fmt.Sprintf("%c - %v", alphabet[i], name), termbox.ColorWhite, termbox.ColorBlack)
	}
}

// Create a new use screen.
func newUseScreen(display display) *screen {
	return &screen{
		display: display,
		panels:  []panel{newUsePanel(display)},
	}
}

// Panel that renders the use list.
type usePanel struct {
	display display
}

// Create a new usePanel.
func newUsePanel(display display) *usePanel {
	return &usePanel{display: display}
}

func (u *usePanel) HandleInput(tboxev termbox.Event) (game.Command, error) {
	if tboxev.Type == termbox.EventKey && tboxev.Key == termbox.KeyEsc {
		return game.ModeCommand{Mode: game.ModeHud}, nil
	} else if ch := tboxev.Ch; ch != 0 {
		opt := selectOption(ch)
		if opt != -1 {
			return game.MenuCommand{Option: opt}, nil
		}
	}
	return nocommand()
}

// Listens to nothing.
func (u *usePanel) HandleEvent(e game.Event) {
}

// Render the menu.
func (m *usePanel) Render(g *game.Game) {
	inv := g.Player.Packer.Inventory()
	renderInventory(m.display, "Use what?", inv)
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
