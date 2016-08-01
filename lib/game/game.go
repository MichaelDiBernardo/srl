package game

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Among other things, a Game serves as a factory for all types of game
// objects.
type ObjFactory interface {
	NewObj(spec *Spec) *Obj
}

// Backend for a single game.
type Game struct {
	Player *Obj
	Level  *Level
	Events *EventQueue
}

// Create a new game.
func NewGame() *Game {
	eq := newEventQueue()

	game := &Game{Events: eq}
	game.Player = game.NewObj(PlayerSpec)

	level := NewLevel(40, 40, game, TestLevel)
	game.Level = level

	level.Place(game.Player, math.Pt(1, 1))
	return game
}

// Create a new object for use in this game.
func (g *Game) NewObj(spec *Spec) *Obj {
	return newObj(spec, g.Events)
}

// Handle a command from the client, and then evolve the world.
func (w *Game) Handle(e Command) {
	switch e {
	case CommandMoveN:
		w.Player.Mover.Move(math.Pt(0, -1))
	case CommandMoveS:
		w.Player.Mover.Move(math.Pt(0, 1))
	case CommandMoveW:
		w.Player.Mover.Move(math.Pt(-1, 0))
	case CommandMoveE:
		w.Player.Mover.Move(math.Pt(1, 0))
	}
	w.Level.Evolve()
}

// A command given _to_ the game.
type Command int

const (
	_ Command = iota
	CommandQuit
	CommandMoveN
	CommandMoveNE
	CommandMoveE
	CommandMoveSE
	CommandMoveS
	CommandMoveSW
	CommandMoveW
	CommandMoveNW
)

// Events are complex objects (unlike commands); you have to type-assert them
// to their concrete types to get at their payloads.
type Event interface{}

// A message that we want to show up in the message console.
type MessageEvent struct {
	Text string
}

// A queue of events that are produced by the game and consumed by the client.
// There are first-class verbs on this queue for each of the event types that
// the game needs to send; nothing pushes directly to the queue.
type EventQueue struct {
	q *list.List
}

// Create a new event queue.
func newEventQueue() *EventQueue {
	return &EventQueue{q: list.New()}
}

// The number of events waiting to be consumed.
func (eq *EventQueue) Len() int {
	return eq.q.Len()
}

// Is this queue empty of events?
func (eq *EventQueue) Empty() bool {
	return eq.Len() == 0
}

// The next event in the queue. Calling this will remove it from the queue and
// furnish the event as a result.
func (eq *EventQueue) Next() Event {
	el := eq.q.Front()
	eq.q.Remove(el)
	return el.Value.(Event)
}

// Send a message to be rendered in the message console.
func (eq *EventQueue) Message(msg string) {
	eq.push(&MessageEvent{Text: msg})
}

// Push a new event onto the queue.
func (eq *EventQueue) push(e Event) {
	eq.q.PushBack(e)
}
