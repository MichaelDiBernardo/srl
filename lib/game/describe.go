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
	melee := fmt.Sprintf("%s%d", extrasign(atk.Melee), atk.Melee)
	dam := ""
	if atk.Damroll != ZeroDice {
		dam = "," + atk.Damroll.String()
	}
	return "(" + melee + dam + ")"
}

func (def Defense) Describe() string {
	evasion := fmt.Sprintf("%s%d", extrasign(def.Evasion), def.Evasion)
	prot := ""
	if len(def.ProtDice) > 0 {
		low, high := 0, 0
		for _, dice := range def.ProtDice {
			low += dice.Dice
			high += dice.Dice * dice.Sides
		}
		prot = "," + fmt.Sprintf("%d-%d", low, high)
	}
	return "[" + evasion + prot + "]"
}

func extrasign(x int) string {
	if x >= 0 {
		return "+"
	}
	return ""
}