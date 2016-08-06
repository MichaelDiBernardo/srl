package console

import (
	"github.com/MichaelDiBernardo/srl/lib/client"
)

// Screens have a similar API to Clients, because the console client defers to
// the current screen in order to get stuff done.
type screen interface {
	client.Renderer
	client.Commander
	client.EventHandler
}

// Panels do everything a client does except generate commands from the user;
// that's the job of the screen, which is composed of panels.
type panel interface {
	client.Renderer
	client.EventHandler
}

var alphabet = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
