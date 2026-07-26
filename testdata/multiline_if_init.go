package fixtures

import "time"

type MempoolTx struct {
	id      [32]byte
	expires time.Time
	tx      any
}

func insert(tx *MempoolTx) error { return nil }

// Single-line init: OK.
func singleLine() {
	if err := insert(nil); err != nil {
		panic(err)
	}
}

// Multi-return single-line init: OK.
func multiReturnSingleLine() {
	if x, err := 1, insert(nil); err != nil {
		_ = x
		panic(err)
	}
}

// Multi-line init with struct literal.
func multiLineStructLiteral() {
	if err := insert(&MempoolTx{ // MATCH /if-init statement should not span multiple lines/
		expires: time.Now().Add(time.Hour),
	}); err != nil {
		panic(err)
	}
}

// Multi-line init with long call chain.
func multiLineCallChain() {
	if err := insert( // MATCH /if-init statement should not span multiple lines/
		nil,
	); err != nil {
		panic(err)
	}
}

// Nested: outer is single-line (OK), inner is multi-line.
func nested() {
	if err := insert(nil); err != nil {
		if err2 := insert(&MempoolTx{ // MATCH /if-init statement should not span multiple lines/
			expires: time.Now(),
		}); err2 != nil {
			panic(err2)
		}
	}
}

// No init clause: OK.
func noInit() {
	x := true
	if x {
		_ = x
	}
}

// Init is a simple assignment (single line): OK.
func simpleAssign() {
	if x := 42; x > 0 {
		_ = x
	}
}
