package game

import "log"

// We expect a speed 2 actor to fully recover in 100 turns.
const RegenTurns = 100

// A thing that keeps track of time passing, and modifies the actor
// accordingly.
type Ticker interface {
	Objgetter
	// Get ready to tick. This is mostly for the player, as their ticker will
	// have to work across levels.
	Init()
	// The absolute delay that has passed on the level.
	Tick(delay int)
}

type ActorTicker struct {
	Trait
	// The last delay value we saw.
	last int
	// regen counter
	regenc int
}

func NewActorTicker(obj *Obj) Ticker {
	return &ActorTicker{Trait: Trait{obj: obj}}
}

func (t *ActorTicker) Init() {
	t.last = 0
}

func (t *ActorTicker) Tick(delay int) {
	diff := delay - t.last
	t.regen(diff)
	t.last = delay
}

func (t *ActorTicker) regen(diff int) {
	t.regenc += diff
	sheet := t.obj.Sheet
	periodt := RegenTurns / sheet.Regen()
	delayPerHp := periodt * getdelay(2) / sheet.MaxHP()
	heal := t.regenc / delayPerHp

	if heal > 0 {
		log.Printf("id%d. I healed %d", t.obj.id, heal)
		sheet.Heal(heal)
		t.regenc -= heal * delayPerHp
	}
}
