package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"testing"
)

// Creates a new game that especially useful for testing; makes a 4x4 map with
// a wall border, and places the player at 2,2.
func newTestGame() *Game {
	return newTestGameWith(SquareLevel)
}

func newTestGameWith(func(*Level) *Level) *Game {
	g := NewGame()
	g.mode = ModeHud
	// Manully do the Start() stuff so we can pick the level.
	g.Player = g.NewObj(PlayerSpec)
	g.Level = NewLevel(4, 4, g, SquareLevel)
	return g
}

// Game tests
func TestKillingPlayerEndsGame(t *testing.T) {
	g := newTestGame()
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

	g := newTestGame()

	obj := g.NewObj(gtActorSpec)

	pos := math.Pt(1, 1)
	g.Level.Place(obj, pos)

	g.Kill(obj)

	if g.Level.At(pos).Actor != nil {
		t.Error(`Actor's previous tile had tile.Actor != nil`)
	}

	// Only actor left should be the player.
	if len(g.Level.actors) > 1 {
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
