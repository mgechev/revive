package fixtures

import (
	"fmt"
	"time"
)

var c chan int

func timeAfterLeak() {
	select {
	case m := <-c:
		handle(m)
	case <-time.After(10 * time.Second): // shall not match /the underlying goroutine of time.After() is not garbage-collected until timer expiration, prefer NewTimer+Timer.Stop/
		fmt.Println("timed out")
	}
}
