// Test data for the early-return rule with allowJump option enabled

package fixtures

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func fn1() {
	if cond { //MATCH /if c { ... } can be rewritten if !c { return } ... to reduce nesting/
		println()
		println()
		println()
	}
}

func fn2() {
	for {
		if cond { //MATCH /if c { ... } can be rewritten if !c { continue } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
}

func fn3() {
	for {
		// can't flip cond2 here because the cond1 branch would flow into it
		if cond1 {
			println()
		} else if cond2 {
			println()
			println()
			println()
		}
	}
}

func fn4() {
	for {
		// cond1 branch continues here so this is ok
		if cond1 {
			println()
			continue
		} else if cond2 { //MATCH /if c { ... } can be rewritten if !c { continue } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
}

func fn5() {
	for {
		// no point flipping cond here we only unnest one statement and need to introduce one new nested statement (continue) to do it
		if cond {
			println()
		}
	}
}

func fn6() {
	for {
		if x, ok := foo(); ok { //MATCH /if c { ... } can be rewritten if !c { continue } ... to reduce nesting (move short variable declaration to its own line if necessary)/
			println(x)
			println(x)
			println(x)
		}
	}
}

func fn7() {
	for i := 0; i < 10; i++ {
		if cond { //MATCH /if c { ... } can be rewritten if !c { continue } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
}

func fn8() {
	for range c {
		if cond { //MATCH /if c { ... } can be rewritten if !c { continue } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
}

func fn9() {
	fn := func() {
		if cond { //MATCH /if c { ... } can be rewritten if !c { return } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
	fn()
}

func fn10() {
	switch {
	case cond:
		if foo() { //MATCH /if c { ... } can be rewritten if !c { break } ... to reduce nesting/
			println()
			println()
			println()
		}
	default:
		if bar() { //MATCH /if c { ... } can be rewritten if !c { break } ... to reduce nesting/
			println()
			println()
			println()
		}
	}
}

func fn11() {
	if a() {
		println()
		os.Exit(1)
	}
}

func fn12() {
	if a() {
		println()
		return
	}
}

func fn13() {
	if err := a(); err != nil {
		println()
		panic(err)
	}
}

func fn14() {
	if err := a(); err != nil {
		println()
		log.Fatal(err)
	}
}

func fn15() {
	if err := a(); err != nil {
		println()
		log.Panic(err)
	}
}

func fn16() {
	if err := a(); err != nil { //MATCH /if c { ... } can be rewritten if !c { return } ... to reduce nesting (move short variable declaration to its own line if necessary)/
		println()
		println()
		log.Panic(err)
	}
}

func fn17() {
	if err := a(); err != nil { //MATCH /if c { ... } can be rewritten if !c { return } ... to reduce nesting (move short variable declaration to its own line if necessary)/
		println()
		println()
		println()
		panic(err)
	}
}

func MustEncode[T any](w http.ResponseWriter, status int, v T) {
	if err := Encode(w, status, v); err != nil {
		slog.Error("Error encoding response", "err", err)
		return
	}
}

func (c *client) renewAuthInfo() {
	err := RenewLease(c.ctx, c, "auth", c.authCreds, func() (*hashiVault.Secret, error) {
		authInfo, err := c.auth(c.v)
		if err != nil {
			return nil, fmt.Errorf("unable to renew auth info: %w", err)
		}

		c.authCreds = authInfo

		return authInfo, nil
	})
	if err != nil {
		slog.Error("unable to renew auth info", slog.String(loggingKeyError, err.Error()))
		os.Exit(1)
	}
}
