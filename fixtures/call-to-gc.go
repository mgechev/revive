package fixtures

import (
	"fmt"
	"runtime"
)

func GC() {
}

func foo() {
	fmt.Println("just testing")
	GC()
	fixtures.GC()
	runtime.Goexit()
	runtime.GC() // MATCH /explicit call to the garbage collector/
}
