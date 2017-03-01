package game

import (
	"strings"
)

// Return a possessive for the given noun.
func poss(noun string) string {
	if strings.HasSuffix(noun, "s") {
		return noun + "'"
	} else {
		return noun + "'s"
	}
}
