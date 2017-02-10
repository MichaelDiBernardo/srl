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
	// The effects currently active on this actor.
	Effects map[Effect]*ActiveEffect
}

func NewActorTicker(obj *Obj) Ticker {
	t := &ActorTicker{
		Trait:   Trait{obj: obj},
		Effects: map[Effect]*ActiveEffect{},
	}
	t.AddEffect(EffectBaseRegen, 0)
	return t
}

func (t *ActorTicker) Tick(delay int) {
	// Non-time-related things.
	if seer := t.obj.Senser; seer != nil {
		seer.CalcFields()
	}
	// Time-related things.
	if delay < t.last {
		// We've been placed into a new level.
		t.last = 0
	}
	diff := delay - t.last

	for _, ae := range t.Effects {
		ae.OnTick(ae, t, diff)
	}
	t.last = delay
}

func (t *ActorTicker) AddEffect(e Effect, counter int) {
	// TODO: Active effects should be cumulative.
	// TODO: Call OnBegin.
	t.Effects[e] = NewActiveEffect(e, counter)
}
