package game

// We expect a speed 2 actor to fully recover in 100 turns.
const RegenPeriod = 100

// A thing that keeps track of time passing, and modifies the actor
// accordingly.
type Ticker interface {
	Objgetter
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

func (t *ActorTicker) Tick(delay int) {
	// We've been placed into a new level.
	if delay < t.last {
		t.last = 0
	}
	diff := delay - t.last
	t.regen(diff)
	t.last = delay
}

func (t *ActorTicker) regen(diff int) {
	sheet := t.obj.Sheet
	regen := sheet.Regen()

	t.regenc += regen * diff
	delayPerHp := RegenPeriod * GetDelay(2) / sheet.MaxHP()
	heal := t.regenc / delayPerHp

	if heal > 0 {
		sheet.Heal(heal)
		t.regenc -= heal * delayPerHp
	}
}
