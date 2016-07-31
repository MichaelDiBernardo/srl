package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
	"math/rand"
)

// A thing that can move given a specific direction.
type Mover interface {
	Move(dir math.Point) bool
}

// A dummy mover used in cases where a thing can't move.
type nullMover struct {
}

// Do nothing and return false.
func (_ *nullMover) Move(dir math.Point) bool {
	return false
}

// Singleton instance of the null mover.
var theNullMover = &nullMover{}

// Constructor for null movers.
func NewNullMover(_ *Obj) Mover {
	return theNullMover
}

// A universally-applicable mover for actors.
type ActorMover struct {
	obj *Obj
}

// Constructor for actor movers.
func NewActorMover(obj *Obj) Mover {
	return &ActorMover{obj: obj}
}

// Try to move the actor. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	moved := obj.Level.Place(obj, endpos)
	if !moved {
		p.obj.Events.Message(fmt.Sprintf("%s says 'ow'.", p.obj.Spec.Name))
	}
	return moved
}

// A thing that can move given a specific direction.
type AI interface {
	Act(l *Level) bool
}

// A dummy AI used when a thing doesn't need a computer to think for it.
type nullAI struct {
}

// Do nothing and return false.
func (_ *nullAI) Act(l *Level) bool {
	return false
}

// Singleton instance of the null mover.
var theNullAI = &nullAI{}

// Constructor for null movers.
func NewNullAI(_ *Obj) AI {
	return theNullAI
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	obj *Obj
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{obj: obj}
}

// Move in any of the 8 directions with uniform chance. Does not take walls
// etc. in account so this will happily try to bump into things.
func (ai *RandomAI) Act(l *Level) bool {
	x, y := rand.Intn(3)-1, rand.Intn(3)-1
	dir := math.Pt(x, y)
	if dir == math.Origin {
		return ai.Act(l)
	}
	log.Printf("AI: Moving from %v by %v", ai.obj.Pos(), dir)

	return ai.obj.Mover.Move(dir)
}

// Accessors for an actor's stats.
type Stats interface {
	Str() int
	Agi() int
	Vit() int
	Mnd() int
}

// Single implementation of this for now; will probably have separate
// implementations for monsters and player when things get more complicated.
type stats struct {
	str int
	agi int
	vit int
	mnd int
	obj *Obj
}

// Given a copy of a stats literal, this will return a function that will bind
// the owner of the stats to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewActorStats(stats stats) func(*Obj) Stats {
	return func(o *Obj) Stats {
		stats.obj = o
		return &stats
	}
}

func (s *stats) Str() int {
	return s.str
}

func (s *stats) Agi() int {
	return s.agi
}

func (s *stats) Vit() int {
	return s.vit
}

func (s *stats) Mnd() int {
	return s.mnd
}

// Stats to assign if the thing has no stats.
var nullstats = &stats{}

func NewNullStats(_ *Obj) Stats {
	return nullstats
}

// A 'character sheet' for an actor. This is where all attributes derived from
// stats + equipment are stored.
type Sheet interface {
	Melee() int
	Evasion() int
	HP() int
	MaxHP() int
	MP() int
	MaxMP() int
}

type sheet struct {
	obj *Obj
	hp  int
	mp  int
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet sheet

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{obj: obj}
	ps.hp = ps.MaxHP()
	ps.mp = ps.MaxMP()
	return ps
}

func (p *PlayerSheet) Melee() int {
	return p.obj.Stats.Agi()
}

func (p *PlayerSheet) Evasion() int {
	return p.obj.Stats.Agi()
}

func (p *PlayerSheet) HP() int {
	return p.hp
}

func (p *PlayerSheet) MP() int {
	return p.mp
}

func (p *PlayerSheet) MaxHP() int {
	return 10 * (1 + p.obj.Stats.Vit())
}

func (p *PlayerSheet) MaxMP() int {
	return 10 * (1 + p.obj.Stats.Mnd())
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	sheet
	melee   int
	evasion int
	maxhp   int
	maxmp   int
}

// Given a copy of a MonsterSheet literal, this will return a function that will bind
// the owner of the sheet to it at object creation time. See the syntax for
// this in actor_spec.go.
func NewMonsterSheet(sheet MonsterSheet) func(*Obj) Sheet {
	return func(o *Obj) Sheet {
		sheet.obj = o
		sheet.hp = sheet.maxhp
		sheet.mp = sheet.maxmp
		return &sheet
	}
}

func (m *MonsterSheet) Melee() int {
	return m.melee
}

func (m *MonsterSheet) Evasion() int {
	return m.evasion
}

func (m *MonsterSheet) HP() int {
	return m.hp
}

func (m *MonsterSheet) MP() int {
	return m.mp
}

func (m *MonsterSheet) MaxHP() int {
	return m.maxhp
}

func (m *MonsterSheet) MaxMP() int {
	return m.maxmp
}

// Sheet to assign if thing has no sheet.
var nullsheet = &MonsterSheet{}

func NewNullSheet(obj *Obj) Sheet {
	return nullsheet
}
