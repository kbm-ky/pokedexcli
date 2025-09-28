package main

import (
	"math/rand/v2"
)

func attemptCatch(user_exp, target_exp int) bool {
	//some made up probablity function
	scaledDelta := (target_exp - user_exp) * 100 / user_exp
	escapePct := 0
	if scaledDelta < -95 {
		escapePct = 5
	} else if scaledDelta > 95 {
		escapePct = 95
	} else {
		escapePct = 45*scaledDelta/95 + 50
	}

	// log.Printf("Escape chance %d%%\n", escapePct)

	// dice roll
	roll := rand.IntN(100)

	return roll >= escapePct
}
