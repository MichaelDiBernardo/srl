package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
	"math/rand"
)

var OSTMonster ObjSubtype = "Monster"
var OSTPlayer ObjSubtype = "Player"

// A thing that can move given a specific direction.
type Mover interface {
	Objgetter
	Move(dir math.Point) bool
}

// A dummy mover used in cases where a thing can't move.
type nullMover struct {
	Trait
}

// Do nothing and return false.
func (_ *nullMover) Move(dir math.Point) bool {
	return false
}

// Constructor for null movers.
func NewNullMover(obj *Obj) Mover {
	return &nullMover{Trait: Trait{obj: obj}}
}

// A universally-applicable mover for actors.
type ActorMover struct {
	Trait
}

// Constructor for actor movers.
func NewActorMover(obj *Obj) Mover {
	return &ActorMover{Trait: Trait{obj: obj}}
}

// Try to move the actor. Return false if the player couldn't move.
func (p *ActorMover) Move(dir math.Point) bool {
	obj := p.obj
	beginpos := obj.Pos()
	endpos := beginpos.Add(dir)

	if !endpos.In(obj.Level) {
		return false
	}

	endtile := obj.Level.At(endpos)
	if other := endtile.Actor; other != nil {
		if opposing := obj.isPlayer() != other.isPlayer(); opposing {
			p.obj.Fighter.Hit(other.Fighter)
		}
		return false
	}

	moved := obj.Level.Place(obj, endpos)
	return moved
}

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Act(l *Level) bool
}

// A dummy AI used when a thing doesn't need a computer to think for it.
type nullAI struct {
	Trait
}

// Do nothing and return false.
func (_ *nullAI) Act(l *Level) bool {
	return false
}

// Constructor for null movers.
func NewNullAI(obj *Obj) AI {
	return &nullAI{Trait: Trait{obj: obj}}
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	Trait
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{Trait: Trait{obj: obj}}
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
	Objgetter
	Str() int
	Agi() int
	Vit() int
	Mnd() int
}

// Single implementation of this for now; will probably have separate
// implementations for monsters and player when things get more complicated.
type stats struct {
	Trait
	str int
	agi int
	vit int
	mnd int
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

func NewNullStats(obj *Obj) Stats {
	return &stats{Trait: Trait{obj: obj}}
}

// A 'character sheet' for an actor. This is where all attributes derived from
// stats + equipment are stored.
type Sheet interface {
	Objgetter
	Melee() int
	Evasion() int
	HP() int
	MaxHP() int
	MP() int
	MaxMP() int

	Hurt(dmg int)
}

type sheet struct {
	Trait
	hp int
	mp int
}

// Sheet used for player, which has a lot of derived attributes.
type PlayerSheet sheet

func NewPlayerSheet(obj *Obj) Sheet {
	ps := &PlayerSheet{Trait: Trait{obj: obj}}
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

func (p *PlayerSheet) Hurt(dmg int) {
	p.hp -= dmg
}

// Sheet used for monsters, which have a lot of hardcoded attributes.
type MonsterSheet struct {
	Trait
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

func (m *MonsterSheet) Hurt(dmg int) {
	m.hp -= dmg
}

func NewNullSheet(obj *Obj) Sheet {
	return &MonsterSheet{Trait: Trait{obj: obj}}
}

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
}

// An attacker that works for all actors.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{Trait: Trait{obj: obj}}
}

func (a *ActorFighter) Hit(other Fighter) {
	dmg := 1
	other.Obj().Sheet.Hurt(dmg)
	msg := fmt.Sprintf("%v hit %v (%d).", a.obj.Spec.Name, other.Obj().Spec.Name, dmg)
	a.obj.Events.Message(msg)
}

type NullFighter struct {
	Trait
}

func NewNullFighter(obj *Obj) Fighter {
	return &NullFighter{Trait: Trait{obj: obj}}
}

func (n *NullFighter) Hit(other Fighter) {
}
