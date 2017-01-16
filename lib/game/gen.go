package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
)

// Given the floor # of the level and the "wiggle" (i.e. permissible range
// outside the floor), this will generate n "groups" of objects made from specs
// that are guaranteed not to be outside [floor - wiggle, floor + wiggle] based
// on its given floor. Group sizes are taken from the GroupSize entry for each
// spec -- if this is a monster, it's intended to be the pack size, and if it's
// an item it is intended to be the stack size.
func Generate(n, floor, wiggle int, specs []*Spec, g *Game) [][]*Obj {
	low, high := floor-wiggle, floor+wiggle
	log.Printf("Generate: %d groups, %d specs, floors %d-%d", n, len(specs), low, high)
	candidates := make([]*Spec, 0)

	for _, spec := range specs {
		if spec.Gen.Findable(low, high) {
			candidates = append(candidates, spec)
		}
	}

	ncandidates := len(candidates)
	generated := make([][]*Obj, 0)

	log.Printf("Generate: %d candidates at floor %d", ncandidates, floor)

	// RIP :(
	if ncandidates == 0 {
		log.Print("Generate: No candidates! Returning no groups.")
		return generated
	}

	for i := 0; i < n; i++ {
		selected := candidates[RandInt(0, ncandidates)]
		gsize := math.Max(1, selected.Gen.GroupSize)
		group := make([]*Obj, 0, gsize)

		for j := 0; j < gsize; j++ {
			group = append(group, g.NewObj(selected))
		}
		generated = append(generated, group)
	}
	log.Print("Generate: done.")

	return generated
}
