package console

import (
	"github.com/nsf/termbox-go"
)

// A thin wrapper interface over termbox methods and primitive methods we've
// built over it. It is not an attempt to completely abstract away the
// display's backing implementation. It is mostly so that we can inject fake
// displays in tests, without worrying about accidentally loading the termbox
// environment.
type display interface {
	Init() error
	Close()
	Flush() error
	PollEvent() termbox.Event
	Clear(fg, bg termbox.Attribute) error
	SetCell(x, y int, ch rune, fg, bg termbox.Attribute)
	Write(x, y int, text string, fg, bg termbox.Attribute)
}

// A termbox-backed display.
type tbdisplay struct {
}

// Wrapper over SetCell, the primary termbox drawing method.
func (d *tbdisplay) SetCell(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
}

// Write a string to the console.
func (d *tbdisplay) Write(x, y int, text string, fg, bg termbox.Attribute) {
	i := 0
	for _, r := range text {
		d.SetCell(x+i, y, r, fg, bg)
		i++
	}
}

func (d *tbdisplay) Clear(fg, bg termbox.Attribute) error {
	return termbox.Clear(fg, bg)
}

// Init the display.
func (d *tbdisplay) Init() error {
	return termbox.Init()
}

// Teardown the display.
func (d *tbdisplay) Close() {
	termbox.Close()
}

// Force-render the display.
func (d *tbdisplay) Flush() error {
	return termbox.Flush()
}

// Pull an interactive event from the display.
func (d *tbdisplay) PollEvent() termbox.Event {
	return termbox.PollEvent()
}

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
