package console

import (
	"testing"
)

func TestMessagePanelHasLimit(t *testing.T) {
	size := 2
	sut := newMessagePanel(size)
	sut.Message("hi")
	sut.Message("bye")
	sut.Message("foo")

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
