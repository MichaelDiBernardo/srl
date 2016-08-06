package console

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
	"testing"
)

func TestMessagePanelHandleMessageEvent(t *testing.T) {
	sut := newMessagePanel(1, &fakedisplay{})
	ev := game.MessageEvent{Text: "hi"}
	sut.Handle(ev)

	if l := sut.lines.Len(); l != 1 {
		t.Errorf(`MessageEvent added %d lines, want 1`, l)
	}

	if s := sut.lines.Front().Value.(string); s != "hi" {
		t.Errorf(`Bottom line was %v, want 'hi'`, s)
	}
}

func TestMessagePanelHasLimit(t *testing.T) {
	size := 2
	sut := newMessagePanel(size, &fakedisplay{})
	sut.message("hi")
	sut.message("bye")
	sut.message("foo")

	if l := sut.lines.Len(); l != size {
		t.Errorf(`MessagePanel held %d messages, want %d`, l, size)
	}

	if s := sut.lines.Front().Value.(string); s != "bye" {
		t.Errorf(`Bottom line was %v, want 'bye'`, s)
	}

	if s := sut.lines.Back().Value.(string); s != "foo" {
		t.Errorf(`Bottom line was %v, want 'foo'`, s)
	}
}
