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

// Do I have anything in this collection of resists that will resist 'effect'?
func (effects Effects) Resists(effect Effect) bool {
	resist := ResistMap[effect]
	return effects.Has(resist)
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
