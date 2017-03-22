package game

// Does all of the required upkeep to an actor before they take their turn.
type Ticker interface {
	Objgetter
	// Notify the actor that 'delay' time has passed. This must be called
	// exactly once per actor turn, before the actor acts.
	Tick(delay int)
	// Add an active effect to this actor for a duration of 'counter'. It is up
	// to the corresponding ActiveEffect to decide how to interpret the units
	// of counter (may be in delay, may be in turns, etc.)
	AddEffect(e Effect, counter int)
	// Forcibly remove an activeeffect, if it exists. Returns 'true' if there
	// was a nonzero counter to remove.
	RemoveEffect(e Effect) bool
	// What is the counter remaining for the current effect?
	Counter(e Effect) int
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

	// Non-time-related things. This is placed after effects-handling so that
	// we get the most up-to-date field calculation (effects can alter sight,
	// like blindness or see-invis.)
	if seer := t.obj.Senser; seer != nil {
		seer.CalcFields()
	}

	t.last = delay
}

// Adds a new active effect to this actor.
func (t *ActorTicker) AddEffect(e Effect, counter int) {
	ae, prev := t.Effects[e], 0
	if ae == nil {
		ae = NewActiveEffect(e, counter)
		t.Effects[e] = ae
	} else {
		prev = ae.Counter
		switch ae.Stacks {
		case AEStackAdd:
			ae.Counter += counter
		case AEStackReplace:
			ae.Counter = counter
		case AEStackIgnore:
			// Don't want to run OnBegin when ignoring.
			return
		}
	}

	ae.OnBegin(ae, t, prev)
}

func (t *ActorTicker) RemoveEffect(e Effect) bool {
	ae := t.Effects[e]
	if ae == nil {
		return false
	}
	ae.OnEnd(ae, t)
	delete(t.Effects, e)
	return true
}

func (t *ActorTicker) Counter(e Effect) int {
	if ae := t.Effects[e]; ae == nil {
		return 0
	} else {
		return ae.Counter
	}
}
