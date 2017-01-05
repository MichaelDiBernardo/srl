package game

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
)

// Backend for a single game.
type Game struct {
	Player *Obj
	Level  *Level
	Events *EventQueue
	Depth  int
	mode   Mode
}

// Create a new game.
func NewGame() *Game {
	return &Game{Events: newEventQueue(), Depth: 1}
}

// Temp convenience method to init the game before playing.
func (g *Game) Start() {
	g.mode = ModeHud
	g.Player = g.NewObj(PlayerSpec)
	g.Level = NewDungeon(g)
}

// Create a new object for use in this game.
func (g *Game) NewObj(spec *Spec) *Obj {
	obj := newObj(spec)
	obj.Game = g
	return obj
}

// Handle a command from the client, and then evolve the world.
func (g *Game) Handle(c Command) {
	log.Printf("Handling command: %#v", c)
	controllers[g.mode](g, c)
}

func (g *Game) SwitchMode(m Mode) {
	g.mode = m
	// Signal to client that yes, we have switched.
	g.Events.SwitchMode(m)
}

func (g *Game) Kill(actor *Obj) {
	if actor.IsPlayer() {
		g.Events.Message("The quest for the TOWER ends...")
		g.Events.More()
		g.SwitchMode(ModeGameOver)
	} else {
		g.Level.Remove(actor)
	}
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

type TryDropCommand struct {
}

type TryEquipCommand struct {
}

type TryRemoveCommand struct {
}

type TryUseCommand struct {
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
	ModeEquip:     equipController,
	ModeUse:       useController,
	ModeRemove:    removeController,
	ModeDrop:      dropController,
	ModeSheet:     sheetController,
}

// Do stuff when player is actually playing the game.
func hudController(g *Game, com Command) {
	evolve := false
	switch c := com.(type) {
	case MoveCommand:
		g.Player.Mover.Move(c.Dir)
		evolve = true
	case TryPickupCommand:
		g.Player.Packer.TryPickup()
	case TryDropCommand:
		g.Player.Packer.TryDrop()
	case TryEquipCommand:
		g.Player.Equipper.TryEquip()
	case TryRemoveCommand:
		g.Player.Equipper.TryRemove()
	case TryUseCommand:
		g.Player.User.TryUse()
	case ModeCommand:
		g.SwitchMode(c.Mode)
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

// Do stuff when player is looking at equipment.
func equipController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		g.Player.Equipper.Equip(c.Option)
	}
}

// Do stuff when player is using an item.
func useController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		g.Player.User.Use(c.Option)
	}
}

// Do stuff when player is looking at body.
func removeController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		// Will probably need to be changed to a SlotCommand or
		// something if slots become complicated enough that
		// type-converting them from ints isn't enough.
		g.Player.Equipper.Remove(Slot(c.Option))
	}
}

// Do stuff when player is trying to drop stuff.
func dropController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		g.Player.Packer.Drop(c.Option)
	}
}

// Do stuff when player is looking at character sheet.
func sheetController(g *Game, com Command) {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
}

// Events are complex objects (unlike commands); you have to type-assert them
// to their concrete types to get at their payloads.
type Event interface{}

// A message that we want to show up in the message console.
type MessageEvent struct {
	Text string
}

// Force the player to --more--.
type MoreEvent struct {
}

// Modes that the game can be in.
type Mode int

const (
	ModeHud Mode = iota
	ModeInventory
	ModePickup
	ModeEquip
	ModeRemove
	ModeDrop
	ModeUse
	ModeSheet
	ModeGameOver
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
	eq.push(MessageEvent{Text: msg})
}

// Tell client we're switching game modes to mode.
func (eq *EventQueue) SwitchMode(mode Mode) {
	eq.push(ModeEvent{Mode: mode})
}

// Tell client to force a --more-- confirm.
func (eq *EventQueue) More() {
	eq.push(MoreEvent{})
}

// Push a new event onto the queue.
func (eq *EventQueue) push(e Event) {
	eq.q.PushBack(e)
}
