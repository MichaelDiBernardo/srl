package game

import (
	"fmt"
	"log"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Act() bool
}

// State-machine-based "AI".
type SMAI struct {
	Trait
	// Fixed attributes for this AI.
	Attribs *Personality
	// My behaviour in the form of a state machine.
	Brain SMAIStateMachine
	// My current state object
	cur smaiStateObj
}

func (s *SMAI) Act() bool {
	t := s.cur.Act(s)
	if t != smaiNoTransition {
		s.transition(t)
	}
	return true
}

func (s *SMAI) transition(trans smaiTransition) {
	nextState, ok := s.Brain[smaiKey{s.cur.State(), trans}]
	if !ok {
		panic(fmt.Sprintf("Bad AI state machine transition: (%v, %v)", s.cur.State(), trans))
	}
	s.cur = newSMAIState(nextState)
	s.cur.Init(s)
}

// Stuff that this AI likes to do.
type Personality struct {
	// How many squares away can I smell things?
	SmellRange int
	// How many squares away can I see things?
	SightRange int
	// How far away will I run from home if I'm territorial?
	ChaseRange int
	// Where is my home if I'm territorial?
	Home math.Point
	// What percent HP do I need to be at before I run away? '25' means '25%'.
	Fear int
}

func NewSMAI(spec SMAI) func(*Obj) AI {
	return func(o *Obj) AI {
		// Copy AI.
		// TODO: Copy personality
		smai := spec
		smai.obj = o
		smai.cur = newSMAIState(smaiUnborn)
		smai.transition(smaiStart)

		return &smai
	}
}

// A state in the machine.
type smaiState int

// A transition from one state to another, which potentially triggers an
// action.
type smaiTransition int

// The "key" used to lookup a start-state:transition pair in the machine.
type smaiKey struct {
	state      smaiState
	transition smaiTransition
}

// A state object that governs how a monster behaves turn-by-turn when in this
// state.
type smaiStateObj interface {
	// Get this state ready to do stuff.
	Init(me *SMAI)

	// Do something for 'me' on this turn. If 'something' means getting me out
	// of this state, the returned transition will be something other than
	// 'smaiNoTransition'
	Act(me *SMAI) smaiTransition

	// Returns the state for this obj.
	State() smaiState
}

// Every state object needs to remember what actual state it represents. We
// embed this into each of them to save us from having to reimplemnt State() on
// every single one.
type smaiSB struct {
	state smaiState
}

func (s smaiSB) State() smaiState {
	return s.state
}

// Create a state object for the given state.
func newSMAIState(state smaiState) smaiStateObj {
	switch state {
	case smaiUnborn:
		return &smaiStateDoNothing{smaiSB: smaiSB{state: smaiUnborn}}
	case smaiStopped:
		return &smaiStateStopped{smaiSB: smaiSB{state: smaiStopped}}
	case smaiWandering:
		return &smaiStateWandering{smaiSB: smaiSB{state: smaiWandering}}
	case smaiChasing:
		return &smaiStateChasing{smaiSB: smaiSB{state: smaiChasing}}
	default:
		panic(fmt.Sprintf("Could not create stateobj for state %v", state))
	}
}

// A specification of a certain behaviour.
type SMAIStateMachine map[smaiKey]smaiState

const (
	// The initial state.
	smaiUnborn smaiState = iota
	// Just chillin after wandering.
	smaiStopped
	// Sitting at my house.
	smaiAtHome
	// Walking to an arbitrary destination.
	smaiWandering
	// Chasing the player.
	smaiChasing
	// Running away from the player.
	smaiFleeing
	// Returning to territorial home.
	smaiGoingHome
)

const (
	// Start doing stuff. This is run when a monster is first spawned.
	smaiStart smaiTransition = iota
	// Pick a new destination and wander towards it.
	smaiWander
	// I found the player!
	smaiFoundPlayer
	// I lost the player!
	smaiLostPlayer
	// I need to get out of here!
	smaiFlee
	// I found the tile I was looking for!
	smaiStopWandering
	// I need to stop running away!
	smaiStopFleeing
	// Dummy transition.
	smaiNoTransition
)

// A very boring state.
type smaiStateDoNothing struct {
	smaiSB
}

func (_ *smaiStateDoNothing) Init(me *SMAI) {}

func (_ *smaiStateDoNothing) Act(me *SMAI) smaiTransition {
	return smaiNoTransition
}

// A state that sits around for a fixed number of turns until wandering again.
type smaiStateStopped struct {
	smaiSB
	// How many more turns should we stay stopped for.
	turns int
}

func (s *smaiStateStopped) Init(me *SMAI) {
	s.turns = RandInt(5, 25)
}

func (s *smaiStateStopped) Act(me *SMAI) smaiTransition {
	if me.obj.Seer.CanSee(me.obj.Game.Player) {
		return smaiFoundPlayer
	}

	s.turns--
	if s.turns <= 0 {
		return smaiWander
	}
	return smaiNoTransition
}

// A state that guides a wandering monster to an arbitrary location.
type smaiStateWandering struct {
	smaiSB
	// The path we're following.
	path Path
	// How many turns we've spent trying to make progress.
	turnsBlocked int
	// Our destination. It's just the last element in path.
	dest math.Point
}

func (s *smaiStateWandering) Init(me *SMAI) {
	// Pick a random tile to wander to.
	level := me.obj.Level
	tile := level.RandomClearTile()

	if tile == nil {
		// We couldn't find a destination.
		s.path = Path{}
		return
	}
	s.findpath(me, tile.Pos)
}

func (s *smaiStateWandering) Act(me *SMAI) smaiTransition {
	if me.obj.Seer.CanSee(me.obj.Game.Player) {
		return smaiFoundPlayer
	}

	// We've reached our destination!
	if len(s.path) == 0 {
		return smaiStopWandering
	}

	// Try to move to our next point.
	mypos, nextpos := me.obj.Pos(), s.path[0]
	if math.ChebyDist(mypos, nextpos) > 1 {
		s.findpath(me, s.dest)
		return s.Act(me)
	}

	dir := nextpos.Sub(mypos)
	err := me.obj.Mover.Move(dir)

	if err == nil {
		s.turnsBlocked = 0
		s.path = s.path[1:]
	} else {
		s.turnsBlocked++
	}

	if s.turnsBlocked > 5 {
		return smaiStopWandering
	}

	return smaiNoTransition
}

func (s *smaiStateWandering) findpath(me *SMAI, dest math.Point) {
	mypos := me.obj.Pos()
	path, ok := me.obj.Level.FindPath(mypos, dest, PathCost)
	if !ok {
		// We can't find our way to our destination. Let's pretend our
		// destination is right here.
		s.path = Path{}
		return
	}
	s.path = path
	s.dest = dest
}

// Chases the player when they're in sight or smell range.
type smaiStateChasing struct {
	smaiSB
	// How many turns in a row we've been chasing without seeing the player.
	turnsUnseen int
}

func (s *smaiStateChasing) Init(me *SMAI) {
	me.obj.Game.Events.Message(fmt.Sprintf("%s shouts!", me.obj.Spec.Name))
	log.Printf("id%d. I see player at %v. Time to chase!!", me.obj.id, me.obj.Game.Player.Pos())
}

func (s *smaiStateChasing) Act(me *SMAI) smaiTransition {
	obj := me.obj
	pos := obj.Pos()

	if obj.Seer.CanSee(obj.Game.Player) {
		dir := math.Origin
		s.turnsUnseen = 0

		// Try chasing by sight.
		playerpos := obj.Game.Player.Pos()

		switch {
		case pos.X < playerpos.X:
			dir.X = 1
		case pos.X > playerpos.X:
			dir.X = -1
		}
		switch {
		case pos.Y < playerpos.Y:
			dir.Y = 1
		case pos.Y > playerpos.Y:
			dir.Y = -1
		}
		log.Printf("id%d. I see player at %v. I'm at %v moving %v", obj.id, playerpos, pos, dir)

		// Try to move. If this doesn't work because we can see the character
		// but the best direction moves us into a wall, we'll just use scent
		// instead. However, that won't count towards the # of turns unseen.
		if err := obj.Mover.Move(dir); err == nil {
			return smaiNoTransition
		} else {
			log.Printf("id%d. I couldn't move %v: %v. Trying smell.", obj.id, err, dir)
		}
	} else {
		s.turnsUnseen++
	}

	// Try chasing by smell.
	maxscent, maxloc, around := 0, math.Origin, me.obj.Level.Around(pos)

	for _, tile := range around {
		if curscent := tile.Flows[FlowScent]; !tile.Feature.Solid && curscent > maxscent {
			maxscent = curscent
			maxloc = tile.Pos
		}
	}

	dir := math.Origin
	if maxscent != 0 {
		dir = maxloc.Sub(pos)
		log.Printf("id%d. I smell player at %v. I'm at %v moving %v", obj.id, maxloc, pos, dir)
	}

	if err := obj.Mover.Move(dir); err != nil {
		log.Printf("id%d. I couldn't move %v: %v", obj.id, dir, err)
	}
	if s.turnsUnseen > 20 {
		return smaiLostPlayer
	}
	return smaiNoTransition
}

// A wandering monster. Randomly picks destinations to walk to, until it
// detects the player.
var SMAIWanderer = SMAIStateMachine{
	{smaiUnborn, smaiStart}:            smaiStopped,
	{smaiStopped, smaiWander}:          smaiWandering,
	{smaiStopped, smaiFoundPlayer}:     smaiChasing,
	{smaiWandering, smaiStopWandering}: smaiStopped,
	{smaiWandering, smaiFoundPlayer}:   smaiChasing,
	{smaiChasing, smaiLostPlayer}:      smaiStopped,
}
