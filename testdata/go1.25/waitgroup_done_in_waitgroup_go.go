package fixtures

import (
	"fmt"
	"sync"
)

func waitGroupDoneInWaitGroupGo() {
	wg := sync.WaitGroup{}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			wg.Done() // MATCH /do not call wg.Done inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
		})
	}

	wg.Wait()
}
