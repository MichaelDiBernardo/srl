package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
	"log"
)

// A thing that can move given a specific direction.
type AI interface {
	Objgetter
	Act(l *Level) bool
}

// An AI that directs an actor to move completely randomly.
type RandomAI struct {
	Trait
}

// Constructor for random AI.
func NewRandomAI(obj *Obj) AI {
	return &RandomAI{Trait: Trait{obj: obj}}
}

// Move in any of the 8 directions with uniform chance. Does not take walls
// etc. in account so this will happily try to bump into things.
func (ai *RandomAI) Act(l *Level) bool {
	obj := ai.obj
	pos := obj.Pos()
	dir := math.Origin

	// Try chasing by sight.
	if obj.Seer != nil && obj.Seer.CanSee(obj.Game.Player) {
		playerpos := obj.Game.Player.Pos()
		log.Printf("I see you! I'm at %v, ur at %v", pos, playerpos)

		switch {
		case pos.X < playerpos.X:
			dir.X = 1
		case pos.X > playerpos.X:
			dir.X = -1
		}
		switch {
		case pos.Y < playerpos.Y:
			dir.Y = 1
		case pos.Y > playerpos.Y:
			dir.Y = -1
		}
		log.Printf("I'm chasing you: %v", dir)
	} else {
		// Try chasing by smell.
		maxscent, maxloc, around := 0, math.Origin, l.Around(pos)

		for _, tile := range around {
			if curscent := tile.Flows[FlowScent]; curscent > maxscent {
				maxscent = curscent
				maxloc = tile.Pos
			}
		}

		if maxscent != 0 {
			dir = maxloc.Sub(pos)
			log.Printf("I smell you! I'm chasing you: %v", dir)
		}
	}
	return ai.obj.Mover.Move(dir)
}
