package game

import (
	"container/list"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

const MaxFloor = 5

// Backend for a single game.
type Game struct {
	Player   *Obj
	Level    *Level
	Events   *EventQueue
	Progress *Progress
	mode     Mode
}

type Progress struct {
	// Floor player is on.
	Floor int
	// Floor player was previously on.
	PrevFloor int
	// The highest floor this player has been on.
	MaxFloor int
	// How many turns have passed.
	Turns int
}

// Change floor. Returns `true` if this is higher than the player has gone
// before.
func (p *Progress) ChangeFloor(dir int) bool {
	p.PrevFloor = p.Floor
	p.Floor += dir

	if p.Floor > p.MaxFloor {
		p.MaxFloor = p.Floor
		return true
	}

	return false
}

// Create a new game.
func NewGame() *Game {
	return &Game{
		Events: newEventQueue(),
		Progress: &Progress{
			Floor:     1,
			PrevFloor: 1,
			MaxFloor:  1,
			Turns:     0,
		},
	}
}

// Temp convenience method to init the game before playing.
func (g *Game) Start() {
	g.mode = ModeHud
	g.Player = g.NewObj(PlayerSpec)
	g.Player.Equipper.Body().Wear(g.NewObj(Items[1]))
	// TODO: We need an InitPlayer
	g.Player.Learner.(*ActorLearner).gainxp(5000)
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
	evolve := controllers[g.mode](g, c)
	if evolve {
		for {
			g.Level.Evolve()
			g.Progress.Turns++
			// If player is para/stone/whatever, keep evolving the game because
			// they can't do anything right now.
			if g.Player.Sheet.CanAct() {
				break
			}
		}
	}
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

// Switch floors on the player.
func (g *Game) ChangeFloor(dir int) {
	isNewMax := g.Progress.ChangeFloor(dir)
	if isNewMax {
		g.Player.Learner.GainXPFloor(g.Progress.Floor)
	}
	g.Level = NewDungeon(g)
}

// A command given _to_ the game.
type Command interface{}

type QuitCommand struct{}

type MoveCommand struct{ Dir math.Point }

type RestCommand struct{}

type TryPickupCommand struct{}

type TryDropCommand struct{}

type TryEquipCommand struct{}

type TryRemoveCommand struct{}

type TryUseCommand struct{}

type TryShootCommand struct{}

type ModeCommand struct{ Mode Mode }

type MenuCommand struct{ Option int }

type AscendCommand struct{}

type DescendCommand struct{}

type StartLearningCommand struct{}

type CancelLearningCommand struct{}

type FinishLearningCommand struct{}

type LearnSkillCommand struct{ Skill SkillName }

type UnlearnSkillCommand struct{ Skill SkillName }

type NoCommand struct{}

// A controller is a function that handles 'command' using 'game', and returns
// true if a turn should pass due to the player's action. (Some player commands
// might result in no time passing, like cancelling out of a menu or trying to
// move through a wall.
type controller func(*Game, Command) bool

// Controller functions that take commands for each given mode and run them on
// the game.
var controllers = map[Mode]controller{
	ModeDrop:      dropController,
	ModeEquip:     equipController,
	ModeHud:       hudController,
	ModeInventory: inventoryController,
	ModePickup:    pickupController,
	ModeRemove:    removeController,
	ModeSheet:     sheetController,
	ModeShoot:     shootController,
	ModeUse:       useController,
}

// Do stuff when player is actually playing the game.
func hudController(g *Game, com Command) bool {
	evolve := false
	switch c := com.(type) {
	case MoveCommand:
		ok, _ := g.Player.Mover.Move(c.Dir)
		evolve = ok
	case RestCommand:
		g.Player.Mover.Rest()
		evolve = true
	case TryPickupCommand:
		evolve = g.Player.Packer.TryPickup()
	case TryDropCommand:
		g.Player.Packer.TryDrop()
	case TryEquipCommand:
		g.Player.Equipper.TryEquip()
	case TryRemoveCommand:
		g.Player.Equipper.TryRemove()
	case TryUseCommand:
		g.Player.User.TryUse()
	case TryShootCommand:
		g.Player.Shooter.TryShoot()
	case AscendCommand:
		g.Player.Mover.Ascend()
		evolve = true
	case DescendCommand:
		g.Player.Mover.Descend()
		evolve = true
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
	return evolve
}

func shootController(g *Game, com Command) bool {
	evolve := false
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
	return evolve
}

// Do stuff when player is looking at inventory.
func inventoryController(g *Game, com Command) bool {
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
	return false
}

// Do stuff when player is looking at ground.
func pickupController(g *Game, com Command) bool {
	evolve := false
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		evolve = g.Player.Packer.Pickup(c.Option)
	}
	return evolve
}

// Do stuff when player is looking at equipment.
func equipController(g *Game, com Command) bool {
	evolve := false
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		evolve = g.Player.Equipper.Equip(c.Option)
	}
	return evolve
}

// Do stuff when player is using an item.
func useController(g *Game, com Command) bool {
	evolve := false
	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		evolve = g.Player.User.Use(c.Option)
	}
	return evolve
}

// Do stuff when player is looking at body.
func removeController(g *Game, com Command) bool {
	evolve := false

	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		// Will probably need to be changed to a SlotCommand or
		// something if slots become complicated enough that
		// type-converting them from ints isn't enough.
		evolve = g.Player.Equipper.Remove(Slot(c.Option))
	}
	return evolve
}

// Do stuff when player is trying to drop stuff.
func dropController(g *Game, com Command) bool {
	evolve := false

	switch c := com.(type) {
	case ModeCommand:
		g.SwitchMode(c.Mode)
	case MenuCommand:
		evolve = g.Player.Packer.Drop(c.Option)
	}
	return evolve
}

// Do stuff when player is looking at character sheet.
func sheetController(g *Game, com Command) bool {
	switch c := com.(type) {
	case StartLearningCommand:
		change, err := g.Player.Learner.BeginLearning()
		if err != nil {
			panic("Wat")
		}
		g.Events.SkillChange(change)
	case CancelLearningCommand:
		err := g.Player.Learner.CancelLearning()
		if err != nil {
			panic("Wat")
		}
	case FinishLearningCommand:
		err := g.Player.Learner.EndLearning()
		if err != nil {
			panic("Wat")
		}
	case LearnSkillCommand:
		change, _ := g.Player.Learner.LearnSkill(c.Skill)
		g.Events.SkillChange(change)
	case UnlearnSkillCommand:
		change, _ := g.Player.Learner.UnlearnSkill(c.Skill)
		g.Events.SkillChange(change)
	case ModeCommand:
		g.SwitchMode(c.Mode)
	}
	return false
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

// Sent while editing skills.
type SkillChangeEvent struct {
	Change *SkillChange
}

// Modes that the game can be in.
type Mode int

const (
	ModeDrop Mode = iota
	ModeEquip
	ModeGameOver
	ModeHud
	ModeInventory
	ModePickup
	ModeRemove
	ModeSheet
	ModeShoot
	ModeUse
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

func (eq *EventQueue) SkillChange(c *SkillChange) {
	eq.push(SkillChangeEvent{Change: c})
}

// Tell client to force a --more-- confirm.
func (eq *EventQueue) More() {
	eq.push(MoreEvent{})
}

// Push a new event onto the queue.
func (eq *EventQueue) push(e Event) {
	eq.q.PushBack(e)
}
