package fixtures

import (
	"sync"
)

func foo(a int, b float32, c char, d sync.WaitGroup) { // MATCH /sync.WaitGroup passed by value, the function will get a copy of the original one/

}

func bar(a, b sync.WaitGroup) { // MATCH /sync.WaitGroup passed by value, the function will get a copy of the original one/

}

func baz(zz sync.WaitGroup) { // MATCH /sync.WaitGroup passed by value, the function will get a copy of the original one/

}

func ok(zz *sync.WaitGroup) {

}
