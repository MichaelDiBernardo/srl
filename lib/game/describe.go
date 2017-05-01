package game

import (
	"fmt"
	"github.com/MichaelDiBernardo/srl/lib/math"
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
	hit := fmt.Sprintf("%s%d", extrasign(atk.Hit), atk.Hit)
	dam := ""
	if atk.Damroll != ZeroDice {
		dam = "," + atk.Damroll.String()
	}
	return "(" + hit + dam + ")"
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
		for _, dice := range def.CorrDice {
			low -= dice.Dice * dice.Sides
			high -= dice.Dice
		}
		prot = "," + fmt.Sprintf("%d-%d", math.Max(low, 0), math.Max(high, 0))
	}
	return "[" + evasion + prot + "]"
}

func extrasign(x int) string {
	if x >= 0 {
		return "+"
	}
	return ""
}
