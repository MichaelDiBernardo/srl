package game

import (
	"fmt"
	"log"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Init()
	Act() bool
}

// State-machine-based "AI".
type SMAI struct {
	Trait
	// Fixed attributes for this AI.
	Personality *Personality
	// My behaviour in the form of a state machine.
	Brain SMAIStateMachine
	// My current state object
	cur smaiStateObj
}

func NewSMAI(spec SMAI) func(*Obj) AI {
	return func(o *Obj) AI {
		// Copy AI.
		smai := spec
		smai.obj = o
		smai.cur = newSMAIState(smaiUnborn)
		smai.Personality = &Personality{}
		*(smai.Personality) = *(spec.Personality)
		return &smai
	}
}

func (s *SMAI) Init() {
	s.transition(smaiStart)
}

func (s *SMAI) Act() bool {
	// Check for flight and recovery conditions.
	percentHP := s.obj.Sheet.HP() * 100 / s.obj.Sheet.MaxHP()
	if percentHP < s.Personality.Fear && s.cur.State() != smaiFleeing {
		s.transition(smaiFlee)
	}
	if percentHP >= s.Personality.Fear && s.cur.State() == smaiFleeing {
		s.transition(smaiStopFleeing)
	}

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
	// The following things are set externally, in configuration.

	// Roughly how many turns will I spend chasing a player by smell, out of
	// LOS? 0 means I'll only chase in LOS.
	Persistence int
	// What percent HP do I need to be at before I run away? '25' means '25%'.
	Fear int

	// These things are set by the state machine itself.
	// Where is my home if I'm territorial?
	home math.Point
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
	case smaiWaiting:
		return &smaiStateWaiting{smaiSB: smaiSB{state: smaiWaiting}}
	case smaiWandering:
		return &smaiStateWandering{smaiSB: smaiSB{state: smaiWandering}}
	case smaiChasing:
		return &smaiStateChasing{smaiSB: smaiSB{state: smaiChasing}}
	case smaiFleeing:
		return &smaiStateFleeing{smaiSB: smaiSB{state: smaiFleeing}}
	case smaiAtHome:
		return &smaiStateAtHome{smaiSB: smaiSB{state: smaiAtHome}}
	case smaiGoingHome:
		return &smaiStateGoingHome{smaiSB: smaiSB{state: smaiGoingHome}}
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
	smaiWaiting
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
	smaiStopWaiting
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
	// I found my house!
	smaiFoundHome
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
type smaiStateWaiting struct {
	smaiSB
	// How many more turns should we wait for.
	turns int
}

func (s *smaiStateWaiting) Init(me *SMAI) {
	s.turns = RandInt(5, 25)
}

func (s *smaiStateWaiting) Act(me *SMAI) smaiTransition {
	if me.obj.Senser.CanSee(me.obj.Game.Player) {
		return smaiFoundPlayer
	}

	s.turns--
	if s.turns <= 0 {
		return smaiStopWaiting
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
	if me.obj.Senser.CanSee(me.obj.Game.Player) {
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

	if err == ErrMoveBlocked {
		s.turnsBlocked++
	} else {
		s.turnsBlocked = 0
		s.path = s.path[1:]
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
	// How many turns are we willing to spend chasing by scent alone?
	motivation int
}

func (s *smaiStateChasing) Init(me *SMAI) {
	me.obj.Game.Events.Message(fmt.Sprintf("%s shouts!", me.obj.Spec.Name))

	// Figure out how long I'll chase by scent.
	persistence := me.Personality.Persistence

	// 0 persistence is a signal that we only want to chase in LOS.
	if persistence == 0 {
		s.motivation = 0
	} else {
		plow, phigh := math.Max(0, persistence-10), persistence+11
		s.motivation = RandInt(plow, phigh)
	}

	log.Printf("id%d. I see player at %v. I'm at %v. Time to chase!!", me.obj.id, me.obj.Game.Player.Pos(), me.obj.Pos())
}

func (s *smaiStateChasing) Act(me *SMAI) smaiTransition {
	obj := me.obj
	pos := obj.Pos()

	if obj.Senser.CanSee(obj.Game.Player) {
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
		err := obj.Mover.Move(dir)

		if err != ErrMoveBlocked {
			return smaiNoTransition
		} else {
			log.Printf("id%d. I couldn't move %v: %v. Trying smell.", obj.id, err, dir)
		}
	} else {
		s.turnsUnseen++
	}

	if s.turnsUnseen > s.motivation {
		log.Printf("id%d. I gave up on smelling. %d/%d", obj.id, s.turnsUnseen, s.motivation)
		return smaiLostPlayer
	}

	// Try chasing by smell.
	maxscent, maxloc, around := 0, math.Origin, me.obj.Level.Around(pos)

	for _, tile := range around {
		if curscent := tile.Scent; !tile.Feature.Solid && curscent > maxscent {
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
	return smaiNoTransition
}

// Runs away from the player.
type smaiStateFleeing struct {
	smaiSB
	turnsBlocked int
	path         Path
	dest         math.Point
	helpless     bool
}

func (s *smaiStateFleeing) Init(me *SMAI) {
	me.obj.Game.Events.Message(fmt.Sprintf("%s flees!", me.obj.Spec.Name))
	log.Printf("id%d. I'm running!! My pos is %v", me.obj.id, me.obj.Game.Player.Pos())
	s.findsafety(me)
}

func (s *smaiStateFleeing) Act(me *SMAI) smaiTransition {
	// If I'm stuck, don't do anything.
	if s.helpless {
		return smaiNoTransition
	}

	// If I arrived, repath.
	if len(s.path) == 0 {
		s.findsafety(me)
		return s.Act(me)
	}

	// If I lost the path, repath.
	mypos, nextpos := me.obj.Pos(), s.path[0]
	if math.ChebyDist(mypos, nextpos) > 1 {
		s.findsafety(me)
		return s.Act(me)
	}

	// Move.
	dir := nextpos.Sub(mypos)
	err := me.obj.Mover.Move(dir)

	if err == ErrMoveBlocked {
		s.turnsBlocked++
	} else {
		s.turnsBlocked = 0
		s.path = s.path[1:]
	}

	// If we're blocked, try a different direction.
	if s.turnsBlocked > 3 {
		s.findsafety(me)
	}

	return smaiNoTransition
}

func (s *smaiStateFleeing) findsafety(me *SMAI) {
	tile := me.obj.Level.RandomClearTile()
	if tile == nil {
		s.helpless = true
		return
	}

	dest := tile.Pos
	path, ok := me.obj.Level.FindPath(me.obj.Pos(), dest, fleecost)

	if !ok {
		s.helpless = true
		return
	}
	s.path = path
	s.dest = dest
}

// A very homey state.
type smaiStateAtHome struct {
	smaiSB
}

func (s *smaiStateAtHome) Init(me *SMAI) {
	if me.Personality.home != math.Origin {
		log.Printf("id%d. I have a home already: %v.", me.obj.id, me.Personality.home)
		return
	}

	me.Personality.home = me.obj.Pos()
	log.Printf("id%d. I have chosen %v as my home.", me.obj.id, me.Personality.home)
}

func (s *smaiStateAtHome) Act(me *SMAI) smaiTransition {
	// We have an uninvited guest.
	if me.obj.Senser.CanSee(me.obj.Game.Player) {
		log.Printf("id%d. There's an intruder in my house!.", me.obj.id)
		return smaiFoundPlayer
	}
	return smaiNoTransition
}

// This is how we get back home when we're lost.
type smaiStateGoingHome struct {
	smaiSB
	path Path
}

func (s *smaiStateGoingHome) Init(me *SMAI) {
	log.Printf("id%d. Things are clear. I'm going home.", me.obj.id)
	s.findhome(me)
}

func (s *smaiStateGoingHome) Act(me *SMAI) smaiTransition {
	if me.obj.Senser.CanSee(me.obj.Game.Player) {
		log.Printf("id%d. Found the player on my way home", me.obj.id)
		return smaiFoundPlayer
	}

	// We're home!
	if len(s.path) == 0 {
		log.Printf("id%d. I have arrived home! I'm at %v, my home is %v.", me.obj.id, me.obj.Pos(), me.Personality.home)
		return smaiFoundHome
	}

	// Try to move to our next point.
	mypos, nextpos := me.obj.Pos(), s.path[0]
	if math.ChebyDist(mypos, nextpos) > 1 {
		log.Printf("id%d. I got knocked off my homepath at %v to %v. Repathing.", me.obj.id, me.obj.Pos(), me.Personality.home)
		s.findhome(me)
		return s.Act(me)
	}

	dir := nextpos.Sub(mypos)
	err := me.obj.Mover.Move(dir)

	if err == nil {
		s.path = s.path[1:]
		log.Printf("id%d. Moving closer to home. %v -> %v dest %v.", me.obj.id, me.obj.Pos(), dir, me.Personality.home)
	} else {
		log.Printf("id%d. Stuck while moving home: %v %v -> %v dest %v.", me.obj.id, err, me.obj.Pos(), dir, me.Personality.home)
	}
	return smaiNoTransition
}

func (s *smaiStateGoingHome) findhome(me *SMAI) {
	mypos := me.obj.Pos()
	path, ok := me.obj.Level.FindPath(mypos, me.Personality.home, PathCost)
	if !ok {
		// We can't find our way to our destination. Let's pretend our
		// destination is right here.
		s.path = Path{}
		return
	}
	s.path = path
}

// Pathfinding cost function to use when we're running away. This is the same
// as the normal one, but it really, really doesn't like running through the
// player. The player is scawy right now.
func fleecost(l *Level, loc math.Point) int {
	if actor := l.At(loc).Actor; actor != nil && actor.IsPlayer() {
		// I really don't want to run through the player unless I have no
		// choice.
		return 200
	}
	return PathCost(l, loc)
}

// A wandering monster. Randomly picks destinations to walk to, until it
// detects the player.
var SMAIWanderer = SMAIStateMachine{
	{smaiUnborn, smaiStart}:            smaiWaiting,
	{smaiWaiting, smaiStopWaiting}:     smaiWandering,
	{smaiWaiting, smaiFoundPlayer}:     smaiChasing,
	{smaiWaiting, smaiFlee}:            smaiFleeing,
	{smaiWandering, smaiStopWandering}: smaiWaiting,
	{smaiWandering, smaiFoundPlayer}:   smaiChasing,
	{smaiWandering, smaiFlee}:          smaiFleeing,
	{smaiChasing, smaiLostPlayer}:      smaiWaiting,
	{smaiChasing, smaiFlee}:            smaiFleeing,
	{smaiFleeing, smaiStopFleeing}:     smaiChasing,
}

// A territorial monster. Guards its home (spawn) square and returns there when
// player is out of LOS.
var SMAITerritorial = SMAIStateMachine{
	{smaiUnborn, smaiStart}:          smaiAtHome,
	{smaiAtHome, smaiFoundPlayer}:    smaiChasing,
	{smaiAtHome, smaiFlee}:           smaiFleeing,
	{smaiChasing, smaiLostPlayer}:    smaiWaiting,
	{smaiChasing, smaiFlee}:          smaiFleeing,
	{smaiWaiting, smaiStopWaiting}:   smaiGoingHome,
	{smaiWaiting, smaiFoundPlayer}:   smaiChasing,
	{smaiWaiting, smaiFlee}:          smaiFleeing,
	{smaiGoingHome, smaiFoundHome}:   smaiAtHome,
	{smaiGoingHome, smaiFoundPlayer}: smaiChasing,
	{smaiGoingHome, smaiFlee}:        smaiFleeing,
	{smaiFleeing, smaiStopFleeing}:   smaiWaiting,
}
