package game

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
	ended := make([]Effect, 0)

	// Apply each active effect.
	for e, ae := range t.Effects {
		done := ae.OnTick(ae, t, diff)
		if done {
			ae.OnEnd(ae, t)
			ended = append(ended, e)
		}
	}

	// Remove any effects that are no longer active.
	for _, e := range ended {
		delete(t.Effects, e)
	}

	t.last = delay
}

// Adds a new active effect to this actor.
func (t *ActorTicker) AddEffect(e Effect, counter int) {
	if ae := t.Effects[e]; ae == nil {
		ae := NewActiveEffect(e, counter)
		ae.OnBegin(ae, t)
		t.Effects[e] = ae
	} else {
		ae.Counter += counter
	}
}
