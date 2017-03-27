package console

import (
	"github.com/MichaelDiBernardo/srl/lib/game"
	"testing"
)

func TestMessagePanelHandleEventMessageEvent(t *testing.T) {
	sut := newMessagePanel(1, &fakedisplay{})
	ev := game.MessageEvent{Text: "hi"}
	sut.HandleEvent(ev)

	if l := sut.lines.Len(); l != 1 {
		t.Errorf(`MessageEvent added %d lines, want 1`, l)
	}

	if s := sut.lines.Front().Value.(*messageLine).text; s != "hi" {
		t.Errorf(`Bottom line was %v, want 'hi'`, s)
	}
}

func TestMessagePanelHasLimit(t *testing.T) {
	size := 2
	sut := newMessagePanel(size, &fakedisplay{})

	sut.HandleEvent(game.MessageEvent{Text: "hi"})
	sut.HandleEvent(game.MessageEvent{Text: "bye"})
	sut.HandleEvent(game.MessageEvent{Text: "foo"})

	if l := sut.lines.Len(); l != size {
		t.Errorf(`MessagePanel held %d messages, want %d`, l, size)
	}

	if s := sut.lines.Front().Value.(*messageLine).text; s != "bye" {
		t.Errorf(`Bottom line was %v, want 'bye'`, s)
	}

	if s := sut.lines.Back().Value.(*messageLine).text; s != "foo" {
		t.Errorf(`Bottom line was %v, want 'foo'`, s)
	}
}
