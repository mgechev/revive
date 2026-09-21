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

// Issue #1029
func verify2(a any) func() {
	return func() {
		fn := a.(func() error)
		if err := fn(); err != nil {
			panic(err)
		}
	}

}

func mainf() {
	defer verify2(func() error { // MATCH /prefer not to defer chains of function calls/
		return nil
	})()
}

// Issue #1528
func issue1528() {
	var fn func() int
	defer func() {
		fn = func() int {
			return 0
		}
	}()

	fn()
}
