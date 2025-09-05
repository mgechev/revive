package fixtures

import (
	"fmt"
	"log"
	"sync"
)

func forbiddenCallInWgGo() {
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
			panic("don't panic here")
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panic("don't panic here")
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panicf("don't panic here")
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panicln("don't panic here")
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
		})
	}

	wg.Wait()
}
