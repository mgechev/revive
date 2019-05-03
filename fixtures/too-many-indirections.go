package fixtures

import (
	"fmt"
	"os"
)

func tmi() {
	tmi().mstr[2].mstr[1].foo().bar() // MATCH /Too many 4 (>2) indirections in expression/
	if a > myMap[key].mySlice[1].foo() {
		fmt.Printf(myMap[key].mySlice[1])
	}
	os.Exit(mapp[2].bar().s[2].foo()) // MATCH /Too many 3 (>2) indirections in expression/
}
