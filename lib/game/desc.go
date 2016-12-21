package game

import (
	"fmt"
)

func (species Species) Describe() string {
	switch species {
	case SpecHuman:
		return "Human"
	default:
		return "???"
	}
}

func (atk Attack) Describe() string {
	return fmt.Sprintf("(%s%d,%dd%d)",
		extrasign(atk.Melee),
		atk.Melee,
		atk.Damroll.Dice,
		atk.Damroll.Sides,
	)
}

func (def Defense) Describe() string {
	low, high := 0, 0
	for _, dice := range def.ProtDice {
		low += dice.Dice
		high += dice.Dice * dice.Sides
	}
	return fmt.Sprintf("[%s%d,%d-%d]",
		extrasign(def.Evasion),
		def.Evasion,
		low,
		high,
	)
}

func extrasign(x int) string {
	if x >= 0 {
		return "+"
	}
	return ""
}
