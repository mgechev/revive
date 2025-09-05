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
			wg.Done() // MATCH /do not call wg.Done inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			defer wg.Done() // MATCH /do not call wg.Done inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			panic("don't panic here") // MATCH /do not call panic inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panic("don't panic here") // MATCH /do not call log.Panic inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panicf("don't panic here") // MATCH /do not call log.Panicf inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
			log.Panicln("don't panic here") // MATCH /do not call log.Panicln inside wg.Go/
		})
	}

	for i := 1; i <= 5; i++ {
		wg.Go(func() {
			fmt.Println(i)
		})
	}

	wg.Wait()
}
