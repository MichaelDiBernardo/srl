package game

// Rolls a non-combat skillcheck of 'skill' vs 'difficulty. 'challenger' is
// assumed to be initiating the check, and 'defender' is the actor opposing
// with skill 'difficulty'. 'defender' can be null if there is no actor
// opposing the action -- skill checks can be made against circumstances and
// inanimate objects (e.g. perception must beat fixed skill N to unlock a
// door.) Returns won=true if the check succeeded, false otherise. 'by' is the
// residual amount that the check won or lost by. The challenger must roll a
// total score higher than 'difficulty'; a tie results in a loss.
func skillcheck(skill, difficulty int, challenger, defender *Obj) (won bool, by int) {
	by = (skill + DieRoll(1, 10)) - (difficulty + DieRoll(1, 10))
	won = by > 0
	return won, by
}

// Calculates how much a skill roll should be modified because of
// resistances/vulnerabilities to the effect being tested.
func resistmod(num int) int {
	return num * 10
}
