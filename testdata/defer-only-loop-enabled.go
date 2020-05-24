package fixtures

import "errors"

type tt int

func (t tt) m() {}

func deferrer3() {
	for {
		go func() {
			defer println()
		}()
		defer func() {}() // MATCH /prefer not to defer inside loops/
	}

	defer tt.m()

	defer func() error {
		return errors.New("error")
	}()

	defer recover()

	recover()

	defer deferrer()
}
