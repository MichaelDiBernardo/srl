package game

// We expect a speed 2 actor to fully recover in 100 turns.
const RegenPeriod = 100

// Does all of the required upkeep to an actor before they take their turn.
type Ticker interface {
	Objgetter
	// Notify the actor that 'delay' time has passed.
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
	// Non-time-related things.
	if seer := t.obj.Senser; seer != nil {
		seer.CalcFOV()
	}

	// Time-related things.
	if delay < t.last {
		// We've been placed into a new level.
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
