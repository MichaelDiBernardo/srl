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
	mode   Mode
}

// Create a new game.
func NewGame() *Game {
	return &Game{Events: newEventQueue()}
}

// Temp convenience method to init the game before playing.
func (g *Game) Start() {
	g.mode = ModeHud
	g.Player = g.NewObj(PlayerSpec)

	level := NewLevel(40, 40, g, TestLevel)
	g.Level = level

	level.Place(g.Player, math.Pt(1, 1))
}

// Create a new object for use in this game.
func (g *Game) NewObj(spec *Spec) *Obj {
	obj := newObj(spec)
	obj.Game = g
	return obj
}

// Handle a command from the client, and then evolve the world.
func (g *Game) Handle(c Command) {
	controllers[g.mode](g, c)
}

func (g *Game) SwitchMode(m Mode) {
	g.mode = m
	// Signal to client that yes, we have switched.
	g.Events.SwitchMode(m)
}

// A command given _to_ the game.
type Command interface{}

type QuitCommand struct {
}

type MoveCommand struct {
	Dir math.Point
}

type TryPickupCommand struct {
}

type ModeCommand struct {
	Mode Mode
}

type MenuCommand struct {
	Option int
}

// Controller functions that take commands for each given mode and run them on
// the game.
var controllers = map[Mode]func(*Game, Command){
	ModeHud:       hudController,
	ModeInventory: inventoryController,
	ModePickup:    pickupController,
}

// Do stuff when player is actually playing the game.
func hudController(g *Game, com Command) {
	evolve := true
	switch c := com.(type) {
	case MoveCommand:
		g.Player.Mover.Move(c.Dir)
	case TryPickupCommand:
		g.Player.Packer.TryPickup()
		evolve = false
	case ModeCommand:
		g.SwitchMode(c.Mode)
		evolve = false
	}
	if evolve {
		g.Level.Evolve()
	}
}

// Do stuff when player is looking at inventory.
func inventoryController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
}

// Do stuff when player is looking at ground.
func pickupController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		g.Player.Packer.Pickup(c.Option)
	}
}

// Events are complex objects (unlike commands); you have to type-assert them
// to their concrete types to get at their payloads.
type Event interface{}

// A message that we want to show up in the message console.
type MessageEvent struct {
	Text string
}

// Modes that the game can be in.
type Mode int

const (
	ModeHud       Mode = iota // Playing the game
	ModeInventory             // Looking at inventory.
	ModePickup                // Looking at the items on the floor.
)

// Tells the client that we've switched game 'modes'.
type ModeEvent struct {
	Mode Mode
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

// Tell client we're switching game modes to mode.
func (eq *EventQueue) SwitchMode(mode Mode) {
	eq.push(&ModeEvent{Mode: mode})
}

// Push a new event onto the queue.
func (eq *EventQueue) push(e Event) {
	eq.q.PushBack(e)
}
