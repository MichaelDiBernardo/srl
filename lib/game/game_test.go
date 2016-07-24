package game

import (
	"testing"
)

type testEvent struct {
}

func (te *testEvent) EventType() EventType {
	return 1284901284
}

func TestQueueStartsEmpty(t *testing.T) {
	sut := newEventQueue()
	if !sut.Empty() {
		t.Error(`Expected new EventQueue to be empty.`)
	}
}

func TestPushPopEvent(t *testing.T) {
	sut := newEventQueue()
	e1 := &testEvent{}

	sut.push(e1)

	if sut.Empty() {
		t.Error(`Expected not empty after push(e).`)
	}
	if l := sut.Len(); l != 1 {
		t.Errorf(`After push(e) Len() was %d, want 1`, l)
	}

	e2 := sut.Next().(*testEvent)

	if e1 != e2 {
		t.Errorf(`push(e) != Next(); %p != %p`, e1, e2)
	}
}

func TestMessage(t *testing.T) {
	sut := newEventQueue()
	msg := "OMG!!!"

	sut.Message(msg)
	e := sut.Next().(*MessageEvent)

	if actual := e.Text; actual != msg {
		t.Errorf(`Message(msg): Text was %v, want %v`, actual, msg)
	}
}
