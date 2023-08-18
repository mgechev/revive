package fixtures

import "errors"

type tt int

func (t tt) m() {}

func deferrer() {
	for {
		go func() {
			defer println()
		}()
		defer func() {}() // MATCH /prefer not to defer inside loops/
	}

	defer tt.m() // MATCH /be careful when deferring calls to methods without pointer receiver/

	defer func() error {
		return errors.New("error") // MATCH /return in a defer function has no effect/
	}()

	defer recover() // MATCH /recover must be called inside a deferred function, this is executing recover immediately/

	recover() // MATCH /recover must be called inside a deferred function/

	defer deferrer()

	helper := func(_ interface{}) {}

	defer helper(recover()) // MATCH /recover must be called inside a deferred function, this is executing recover immediately/

	// does not work, but not currently blocked.
	defer helper(func() { recover() })
}

// Issue #863

func verify(fn func() error) {
	if err := fn(); err != nil {
		panic(err)
	}
}

func f() {
	defer verify(func() error {
		return nil
	})
}
