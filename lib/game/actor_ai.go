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
	rule, ok := s.Brain[smaiKey{s.cur.State(), trans}]
	if !ok {
		panic(fmt.Sprintf("Bad AI state machine transition: (%v, %v)", s.cur.State(), trans))
	}
	rule.action(s)
	s.cur = newSMAIState(rule.next)
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

// This is something that should happen when we switch from one state to
// another.
type smaiTransitionAction func(me *SMAI)

// The "key" used to lookup a start-state:transition pair in the machine.
type smaiKey struct {
	state      smaiState
	transition smaiTransition
}

// What action and next state should be triggered when a transition is taken?
type smaiRule struct {
	action smaiTransitionAction
	next   smaiState
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
	default:
		panic(fmt.Sprintf("Could not create stateobj for state %v", state))
	}
}

// A specification of a certain behaviour.
type SMAIStateMachine map[smaiKey]smaiRule

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
	log.Printf("Time to rest! I'm going to wait %d turns.", s.turns)
}

func (s *smaiStateStopped) Act(me *SMAI) smaiTransition {
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
	log.Printf("Time to wander!")
	// Pick a random tile to wander to.
	level := me.obj.Level
	tile := level.RandomClearTile()

	if tile == nil {
		// We couldn't find a destination.
		log.Printf("I couldn't find a good tile to wander to.")
		s.path = Path{}
		return
	}
	s.pathfind(me, tile.Pos)
}

func (s *smaiStateWandering) Act(me *SMAI) smaiTransition {
	// We've reached our destination!
	if len(s.path) == 0 {
		log.Printf("Reached destination %v. Time to stop.", s.dest)
		return smaiStopWandering
	}

	// Try to move to our next point.
	mypos, nextpos := me.obj.Pos(), s.path[0]
	if math.ChebyDist(mypos, nextpos) > 1 {
		log.Printf("I was pushed off course! Repathing to %v.", s.dest)
		s.pathfind(me, s.dest)
		return s.Act(me)
	}

	dir := nextpos.Sub(mypos)
	ok := me.obj.Mover.Move(dir)

	if !ok {
		s.turnsBlocked++
	} else {
		s.turnsBlocked = 0
		s.path = s.path[1:]
	}

	if s.turnsBlocked > 5 {
		log.Printf("Giving up on destination %v, I'm stuck at %v.", s.dest, me.obj.Pos())
		return smaiStopWandering
	}

	return smaiNoTransition
}

func (s *smaiStateWandering) pathfind(me *SMAI, dest math.Point) {
	mypos := me.obj.Pos()
	path, ok := me.obj.Level.FindPath(mypos, dest, PathCost)
	if !ok {
		// We can't find our way to our destination. Let's pretend our
		// destination is right here.
		log.Printf("I couldn't find a path from %v to %v.", mypos, dest)
		s.path = Path{}
		return
	}
	s.path = path
	s.dest = dest
	log.Printf("OK! I'm going to walk from %v to %v in %d steps.", mypos, dest, len(path))
}

// TRANSITION ACTIONS
func noTAction(me *SMAI) {
}

// A wandering monster. Randomly picks destinations to walk to, until it
// detects the player.
var SMAIWanderer = SMAIStateMachine{
	{smaiUnborn, smaiStart}:            {noTAction, smaiStopped},
	{smaiStopped, smaiWander}:          {noTAction, smaiWandering},
	{smaiWandering, smaiStopWandering}: {noTAction, smaiStopped},
}

//// An AI that directs an actor to move completely randomly.
//type RandomAI struct {
//	Trait
//}
//
//// Constructor for random AI.
//func NewRandomAI(obj *Obj) AI {
//	return &RandomAI{Trait: Trait{obj: obj}}
//}
//
//// Move in any of the 8 directions with uniform chance. Does not take walls
//// etc. in account so this will happily try to bump into things.
//func (ai *RandomAI) Act(l *Level) bool {
//	obj := ai.obj
//	pos := obj.Pos()
//	dir := math.Origin
//
//	// Try chasing by sight.
//	if obj.Seer != nil && obj.Seer.CanSee(obj.Game.Player) {
//		playerpos := obj.Game.Player.Pos()
//		log.Printf("I see you! I'm at %v, ur at %v", pos, playerpos)
//
//		switch {
//		case pos.X < playerpos.X:
//			dir.X = 1
//		case pos.X > playerpos.X:
//			dir.X = -1
//		}
//		switch {
//		case pos.Y < playerpos.Y:
//			dir.Y = 1
//		case pos.Y > playerpos.Y:
//			dir.Y = -1
//		}
//		log.Printf("I'm chasing you: %v", dir)
//	} else {
//		// Try chasing by smell.
//		maxscent, maxloc, around := 0, math.Origin, l.Around(pos)
//
//		for _, tile := range around {
//			if curscent := tile.Flows[FlowScent]; curscent > maxscent {
//				maxscent = curscent
//				maxloc = tile.Pos
//			}
//		}
//
//		if maxscent != 0 {
//			dir = maxloc.Sub(pos)
//			log.Printf("I smell you! I'm chasing you: %v", dir)
//		}
//	}
//	return ai.obj.Mover.Move(dir)
//}
