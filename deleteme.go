package main

import "time"

func foo() {

	select {
	case <-time.After(2):
		// do something
	}
}
