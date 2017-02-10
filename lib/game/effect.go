package game

// An effect is something that a monster or a piece of equipment can have. This
// includes, brands, resists, status effects, etc.
type Effect uint

// A collection of effects on a monster, piece of equipment, etc.
type Effects []Effect

// Does a collection of effects have this effect?
func (effects Effects) Has(effect Effect) bool {
	for _, e := range effects {
		if e == effect {
			return true
		}
	}
	return false
}

// Do I have anything in this collection of effects that will resist 'effect'?
func (effects Effects) Resists(effect Effect) bool {
	resist := ResistMap[effect]
	return effects.Has(resist)
}

// Do I have anything in this collection of effects that indicates vulnerability to 'effect'?
func (effects Effects) VulnTo(effect Effect) bool {
	vuln := VulnMap[effect]
	return effects.Has(vuln)
}

// Filters out the brands from this collection of effects.
func (effects Effects) Brands() Effects {
	brands := make(Effects, 0)
	for _, e := range effects {
		if Brands.Has(e) {
			brands = append(brands, e)
		}
	}
	return brands
}

// An instance of an effect that is currently affected an actor. These are
// managed by the actor's ticker.
type ActiveEffect struct {
	// The effect counter. Might be turns, accumulated delay, residual damage, etc.
	Counter int
	// What to do when the effect is first inflicted on an actor.
	OnBegin func(*ActiveEffect, *ActorTicker)
	// Responsible for updating Left given the delay diff, plus enforcing the
	// effect. A return value of 'true' indicates that the effect should be
	// terminated.
	OnTick func(*ActiveEffect, *ActorTicker, int) bool
	// What to do when the effect has run its course.
	OnEnd func(*ActiveEffect, *ActorTicker)
}

func NewActiveEffect(e Effect, counter int) *ActiveEffect {
	ae := &ActiveEffect{}
	*ae = ActiveEffects[e]
	ae.Counter = counter
	return ae
}
