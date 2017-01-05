package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Anything that fights in melee.
type Fighter interface {
	Objgetter
	Hit(other Fighter)
}

func hit(attacker Fighter, defender Fighter) {
	atk, def := attacker.Obj().Sheet.Attack(), defender.Obj().Sheet.Defense()
	residual := DieRoll(1, 20) + atk.Melee - DieRoll(1, 20) + def.Evasion
	aname, dname := attacker.Obj().Spec.Name, defender.Obj().Spec.Name

	if residual > 0 {
		crits := residual / atk.CritDiv
		dmg := math.Max(0, atk.RollDamage(crits)-def.RollProt())

		critstr := ""
		if crits > 0 {
			critstr = fmt.Sprintf(" %dx critical!", crits)
		}
		msg := fmt.Sprintf("%v hit %v (%d).%s", aname, dname, dmg, critstr)
		attacker.Obj().Game.Events.Message(msg)

		defender.Obj().Sheet.Hurt(dmg)
	} else {
		msg := fmt.Sprintf("%v missed %v.", aname, dname)
		attacker.Obj().Game.Events.Message(msg)
	}
}

// Melee combat.
type ActorFighter struct {
	Trait
}

func NewActorFighter(obj *Obj) Fighter {
	return &ActorFighter{
		Trait: Trait{obj: obj},
	}
}

func (f *ActorFighter) Hit(other Fighter) {
	hit(f, other)
}