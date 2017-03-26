package game

import (
	"errors"
	"fmt"

	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Tracks xp gained by monsters and items seen and killed.
type Learner interface {
	// Call this when a specific instance of an item or monster has been seen.
	LearnSight(obj *Obj)
	// Call this when a specific instance of a monster has been seen.
	LearnKill(mon *Obj)
	// Gain XP for seeing a new floor.
	LearnFloor(floor int)
	// How much XP does this actor have?
	XP() int
	// How much have they accumulated in total?
	TotalXP() int

	// Start buying skills. This cannot be called again until CancelLearning or
	// EndLearning has been called.
	BeginLearning() (*SkillChange, error)
	// Buy 1 point of sk. If this costs too much, the existing change will be
	// returned along with ErrNotEnoughXP.  Otherwise, will return the current
	// state of the skillchange being requested. This will spend xp and change
	// skill allocations on the character sheet -- these changes will be
	// reverted if CancelLearning is called. Will return
	// ErrNotLearning as err if BeginLearning() was not called before this.
	LearnSkill(sk SkillName) (*SkillChange, error)
	// Refund 1 point of sk. This is only allowed if some of sk has already
	// been bought -- will return the current change and ErrNoPointsLearned if
	// not. Otherwise, will return the current state of the skillchange being
	// requested. This will refund xp and change skill allocations on the
	// character sheet -- these changes will be reverted if CancelLearning is
	// called. Will return ErrNotLearning as err if BeginLearning() was not
	// called before this.
	UnlearnSkill(sk SkillName) (*SkillChange, error)
	// Cancel the current learning change. This will return the character sheet
	// back to what it was before the changes were proposed. It is safe to call
	// BeginLearning again after this has been called.
	CancelLearning() error
	// Finalize the learning changes. It is safe to call BeginLearning again
	// after this has been called.
	EndLearning() error
}

// Used by the player to track what they have seen, and how much XP they have.
type ActorLearner struct {
	Trait
	seen   map[Species]int
	killed map[Species]int
	xp     int
	change *SkillChange
}

// Stores the state of a player's current request to upgrade their skills by
// spending XP.
type SkillChange struct {
	TotalCost int
	Changes   map[SkillName]SkillChangeItem
}

type SkillChangeItem struct {
	Points int
	Cost   int
}

var (
	ErrAlreadyLearning = errors.New("AlreadyLearning")
	ErrNotLearning     = errors.New("NotLearning")
	ErrNotEnoughXP     = errors.New("NotEnoughXP")
	ErrNoPointsLearned = errors.New("NoPointsLearned")
)

// Implementation of player's learner. Monsters don't need this because they
// don't gain or spend XP.
func NewActorLearner(obj *Obj) Learner {
	return &ActorLearner{
		Trait:  Trait{obj: obj},
		seen:   map[Species]int{},
		killed: map[Species]int{},
	}
}

func (l *ActorLearner) XP() int {
	return l.xp
}

func (l *ActorLearner) TotalXP() int {
	return l.xp
}

func (l *ActorLearner) LearnKill(mon *Obj) {
	if genus := mon.Spec.Genus; genus != GenMonster {
		panic(fmt.Sprintf("Obj *v with genus %v is not monster.", mon, genus))
	}

	s := mon.Spec.Species
	n := l.killed[s]

	xpgain := monxp(mon, n)
	l.xp += xpgain
	l.killed[s]++
}

func (l *ActorLearner) LearnSight(obj *Obj) {
	if obj.Seen {
		return
	}

	s, xp := obj.Spec.Species, 0
	n := l.seen[s]

	switch genus := obj.Spec.Genus; genus {
	case GenMonster:
		xp = monxp(obj, n)
	case GenEquipment, GenConsumable:
		xp = itemxp(obj, n)
	default:
		panic(fmt.Sprintf("Obj *v with genus %v is not xpable on sight.", obj, genus))
	}

	l.xp += xp
	l.seen[s]++
}

func (l *ActorLearner) LearnFloor(floor int) {
	l.xp += floor * 50
}

func (l *ActorLearner) BeginLearning() (*SkillChange, error) {
	if l.change != nil {
		return l.change, ErrAlreadyLearning
	}
	l.change = &SkillChange{Changes: map[SkillName]SkillChangeItem{}}
	return l.change, nil
}

func (l *ActorLearner) LearnSkill(sk SkillName) (*SkillChange, error) {
	if l.change == nil {
		return nil, ErrNotLearning
	}

	currskill := l.obj.Sheet.UnmodSkill(sk)
	schange := l.change.Changes[sk]

	cost := skillxp(currskill + 1)
	newtotal := l.change.TotalCost + cost

	if newtotal > l.XP() {
		return l.change, ErrNotEnoughXP
	}

	schange.Points++
	schange.Cost += cost
	l.change.Changes[sk] = schange
	l.change.TotalCost = newtotal

	// Actually modify the player so that these changes will be fully reflected
	// on the character sheet. If they decide to cancel, we'll revert all
	// changes.
	l.xp -= cost
	l.obj.Sheet.SetSkill(sk, currskill+1)

	return l.change, nil
}

func (l *ActorLearner) UnlearnSkill(sk SkillName) (*SkillChange, error) {
	if l.change == nil {
		return nil, ErrNotLearning
	}

	schange := l.change.Changes[sk]

	if schange.Points == 0 {
		return l.change, ErrNoPointsLearned
	}

	// The current value of the skill will actually have been changed by
	// LearnSkill(), so we just need to check the unmodded skill.
	currskill := l.obj.Sheet.UnmodSkill(sk)
	gain := skillxp(currskill)

	schange.Points -= 1
	schange.Cost -= gain
	l.change.Changes[sk] = schange
	l.change.TotalCost -= gain

	// Actually modify the player so that these changes will be fully reflected
	// on the character sheet. If they decide to cancel, we'll revert all
	// changes.
	l.xp += gain
	l.obj.Sheet.SetSkill(sk, currskill-1)

	return l.change, nil
}

func (l *ActorLearner) CancelLearning() error {
	if l.change == nil {
		return ErrNotLearning
	}

	for sk, change := range l.change.Changes {
		curr := l.obj.Sheet.UnmodSkill(sk)
		l.obj.Sheet.SetSkill(sk, curr-change.Points)
	}
	l.xp += l.change.TotalCost

	l.change = nil
	return nil
}

func (l *ActorLearner) EndLearning() error {
	if l.change == nil {
		return ErrNotLearning
	}
	l.change = nil
	return nil
}

// How much XP should you get for seeing or killing mon for the nth time?
func monxp(mon *Obj, n int) int {
	return xpdecay(xpfor(mon), n)
}

// How much xp should you get for seeing item for the nth time?
func itemxp(item *Obj, n int) int {
	if n > 0 {
		return 0
	}
	return xpfor(item)
}

// Raw xp calculation based on native depth.
func xpfor(obj *Obj) int {
	depth := obj.Spec.Gen.First()
	return depth * 10
}

// Decay XP based on number of times seen / killed.
func xpdecay(xp, n int) int {
	return math.Max(1, xp/(n+1))
}

// How much should the nth point of a skill cost to learn?
func skillxp(n int) int {
	return 100 * n
}
