package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

// Game tests
func TestKillingPlayerEndsGame(t *testing.T) {
	g := NewGame()
	g.Start()
	g.Kill(g.Player)

	if m := g.mode; m != ModeGameOver {
		t.Errorf(`game.Kill(player) changed mode to %v; want %v`, m, ModeGameOver)
	}
}

func TestKillingMonsterRemovesIt(t *testing.T) {
	gtActorSpec := &Spec{
		Family:  FamActor,
		Genus:   GenMonster,
		Species: "TestSpecies",
		Name:    "Hi",
		Traits:  &Traits{},
	}

	g := NewGame()

	obj := g.NewObj(gtActorSpec)
	l := NewLevel(4, 4, g, IdentLevel)
	g.Level = l

	pos := math.Pt(1, 1)
	l.Place(obj, pos)

	g.Kill(obj)

	if l.At(pos).Actor != nil {
		t.Error(`Actor's previous tile had tile.Actor != nil`)
	}
	if len(l.actors) > 0 {
		t.Error(`l.actors was not empty after removal.`)
	}

}

// EventQueue tests
type testEvent struct {
}

func TestQueueStartsEmpty(t *testing.T) {
	sut := newEventQueue()
	if !sut.Empty() {
		t.Error(`Expected new EventQueue to be empty.`)
	}
}

func TestPushPopEvent(t *testing.T) {
	sut := newEventQueue()
	e1 := testEvent{}

	sut.push(e1)

	if sut.Empty() {
		t.Error(`Expected not empty after push(e).`)
	}
	if l := sut.Len(); l != 1 {
		t.Errorf(`After push(e) Len() was %d, want 1`, l)
	}

	e2 := sut.Next().(testEvent)

	if e1 != e2 {
		t.Errorf(`push(e) != Next(); %p != %p`, e1, e2)
	}
}

func TestMessage(t *testing.T) {
	sut := newEventQueue()
	msg := "OMG!!!"

	sut.Message(msg)
	e := sut.Next().(MessageEvent)

	if actual := e.Text; actual != msg {
		t.Errorf(`Message(msg): Text was %v, want %v`, actual, msg)
	}
}
