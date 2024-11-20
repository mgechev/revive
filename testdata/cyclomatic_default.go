package pkg

import "log"

func f(x int) bool { // MATCH /function f has cyclomatic complexity 11 (> max enabled 10)/
	if x > 0 && true || false {
		return true
	} else {
		log.Printf("non-positive x: %d", x)
	}
	switch x {
	case 1:
	case 2:
	case 3:
	case 4:
	default:
	}
	return true || true && true
}
