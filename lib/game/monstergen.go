package game

import (
	"log"
)

// Given the depth of the level and the "wiggle" (i.e. permissible range
// outside the depth), this will generate n "groups" of monsters that contain
// that are guaranteed not to be outside [depth - wiggle, depth + wiggle] based
// on its given depth. Group sizes are taken from the GroupSize entry for each
// monster.
func GenMonsters(n, depth, wiggle int, g *Game) [][]*Obj {
	return genmonsters(n, depth, wiggle, Monsters, g)
}

// This one allows you to specify the list of specs to use as input to make it
// possible to test it on different speclists.
func genmonsters(n, depth, wiggle int, specs []*Spec, g *Game) [][]*Obj {
	low, high := depth-wiggle, depth+wiggle
	log.Printf("genmonsters: %d groups, %d specs, depths %d-%d", n, len(specs), low, high)

	candidates := make([]*Spec, 0)

	log.Print("genmonsters: filtering candidates")
	for _, spec := range specs {
		if spec.Gen.Findable(low, high) {
			log.Printf("\tSelected %v", spec.Name)
			candidates = append(candidates, spec)
		}
	}

	ncandidates := len(candidates)
	generated := make([][]*Obj, 0)

	log.Printf("genmonsters: %d candidates at depth %d", ncandidates, depth)

	// RIP :(
	if ncandidates == 0 {
		log.Print("genmonsters: No candidates! Returning no groups.")
		return generated
	}

	log.Print("genmonsters: Creating groups.")
	for i := 0; i < n; i++ {
		selected := candidates[RandInt(0, ncandidates)]
		gsize := selected.Gen.GroupSize
		group := make([]*Obj, 0, gsize)

		log.Printf("\t%d %v", gsize, selected.Name)
		for j := 0; j < gsize; j++ {
			group = append(group, g.NewObj(selected))
		}
		generated = append(generated, group)
	}
	log.Print("genmonsters: done.")

	return generated
}
