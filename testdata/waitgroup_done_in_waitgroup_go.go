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
			wg.Done()
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			defer wg.Done()
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
		})
	}

	wg.Wait()
}
