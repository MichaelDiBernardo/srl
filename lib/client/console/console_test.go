package console

import (
	"github.com/nsf/termbox-go"
)

// A fake display that does nothing and can be used in tests.
type fakedisplay struct {
}

func (d *fakedisplay) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
}

func (d *fakedisplay) Write(x, y int, text string, fg, bg termbox.Attribute) {
}

func (d *fakedisplay) Clear(fg, bg termbox.Attribute) error {
	return nil
}

func (d *fakedisplay) Init() error {
	return nil
}

func (d *fakedisplay) Close() {
}

func (d *fakedisplay) Flush() error {
	return nil
}

func (d *fakedisplay) PollEvent() termbox.Event {
	return termbox.Event{}
}

func (d *fakedisplay) HideCursor() {
}

func (d *fakedisplay) SetCursor(x, y int) {
}
