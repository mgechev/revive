package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// Defer rule.
func TestDefer(t *testing.T) {
	testRule(t, "defer", &rule.DeferRule{})
}

func TestDeferLoopDisabled(t *testing.T) {
	testRule(t, "defer-loop-disabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: []interface{}{[]interface{}{"return", "recover", "call-chain", "method-call"}},
	})
}

func TestDeferOthersDisabled(t *testing.T) {
	testRule(t, "defer-only-loop-enabled", &rule.DeferRule{}, &lint.RuleConfig{
		Arguments: []interface{}{[]interface{}{"loop"}},
	})
}

func TestDemonstrateIneffectiveDefers(t *testing.T) {
	mustPanic := func() {
		if r := recover(); r == nil {
			t.Fatal("recover should not have suppressed the panic")
		}
	}
	mustNotPanic := func() {
		if r := recover(); r != nil {
			t.Fatal("recover should have suppressed the panic, got:", r)
		}
	}
	t.Run("too-deep", func(t *testing.T) {
		recoverer := func() {
			recover()
		}
		helper := func() {
			recoverer() // this does not work, which is the motivating force behind the "must be deferred" lint
		}
		defer mustPanic()
		func() {
			defer helper()
			panic("should escape")
		}()
	})
	t.Run("dynamic-closure", func(t *testing.T) {
		// rarely seen, but also a bad pattern.
		// Go does not know that `func()` calls `recover()`, so it does not generate the necessary capturing code.
		//
		// this is not currently detected or blocked.
		// doing so precisely is probably not possible, but "arguments cannot contain recover() anywhere" may be good.
		defer mustPanic()
		helper := func(recoverer func()) {
			recoverer()
		}
		func() {
			defer helper(func() { recover() }) // seemingly valid, but does not work
			panic("should escape")
		}()
	})
	t.Run("immediate", func(t *testing.T) {
		t.Run("call", func(t *testing.T) {
			defer mustPanic()
			func() {
				defer recover() // recovers immediately, does not suppress panic
				panic("should escape")
			}()
		})
		t.Run("arg", func(t *testing.T) {
			defer mustPanic()
			func() {
				defer t.Log("nothing recovered:", recover()) // recovers immediately because args are collected immediately, does not suppress panic
				panic("should escape")
			}()
		})
	})
	t.Run("correct-immediate", func(t *testing.T) {
		// constructs I've never seen in practice, but these DO work despite violating the lint.
		// the lint currently claims confidence=1, but if legitimate uses like this are found, it may deserve 0.8.
		t.Run("call", func(t *testing.T) {
			defer mustNotPanic()
			func() {
				defer func() {
					defer recover() // called in defer, works but hides result
				}()
				panic("captured but ignored")
			}()
		})
		t.Run("arg", func(t *testing.T) {
			defer mustNotPanic()
			func() {
				defer func() {
					defer t.Log("successfully recovered:", recover()) // called in defer, works
				}()
				panic("captured")
			}()
		})
	})
}
