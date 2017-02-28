package game

import (
	"github.com/MichaelDiBernardo/srl/lib/math"
)

// Rolls a non-combat skillcheck of 'skill' vs 'difficulty. 'challenger' is
// assumed to be initiating the check, and 'defender' is the actor opposing
// with skill 'difficulty'. 'defender' can be null if there is no actor
// opposing the action -- skill checks can be made against circumstances and
// inanimate objects (e.g. perception must beat fixed skill N to unlock a
// door.) Returns won=true if the check succeeded, false otherise. 'by' is the
// residual amount that the check won or lost by. The challenger must roll a
// total score higher than 'difficulty'; a tie results in a loss.
func skillcheck(skill, difficulty int, resists int, challenger, defender *Obj) (won bool, by int) {
	s := skill + skillroll(challenger)
	d := difficulty + resistmod(resists) + skillroll(defender)

	by = s - d
	won = by > 0
	return won, by
}

// Used when an actor is resisting an external effect that doesn't really have
// a direct actor behind it. This always checks will vs difficulty 10.
func savingthrow(defender *Obj, defeffects Effects, effect Effect) bool {
	score := 10
	difficulty := defender.Sheet.Skill(Chi)
	resists := resistmod(defeffects.Resists(effect))
	won, _ := skillcheck(score, difficulty, resists, nil, defender)
	return won
}

// Depending on the blessed/cursed status flags on 'roller.Sheet', this will
// roll a d10 up to twice and take the best if blessed, and the worst if
// cursed. Setting both flags to true has the same effect as setting both to
// false; only one roll will be made. If roller or roller.Sheet is nil, the die
// will only be rolled once.
func skillroll(roller *Obj) int {
	return sroll(roller, 10)
}

// Same as skillroll, but uses d20s. This should be used for all melee,
// evasion, shooting rolls.
func combatroll(roller *Obj) int {
	return sroll(roller, 20)
}

// Actual implementation of skill + melee rolls.
func sroll(roller *Obj, sides int) int {
	roll := DieRoll(1, sides)

	hassheet := roller != nil && roller.Sheet != nil
	blessed := hassheet && roller.Sheet.Blessed()
	cursed := hassheet && roller.Sheet.Cursed()

	if blessed == cursed {
		return roll
	}
	roll2 := DieRoll(1, sides)

	if blessed {
		return math.Max(roll, roll2)
	} else {
		return math.Min(roll, roll2)
	}
}

// Calculates how much a skill roll should be modified because of
// resistances/vulnerabilities to the effect being tested.
func resistmod(num int) int {
	return num * 10
}
