// Test of cognitive complexity.

// Package pkg ...
package pkg

import (
	"fmt"
	ast "go/ast"
	"log"
)

func f(x int) bool { // MATCH /function f has cognitive complexity 3 (> max enabled 0)/
	if x > 0 && true || false {
		return true
	} else {
		log.Printf("non-positive x: %d", x)
	}
	return false
}

func g(f func() bool) string { // MATCH /function g has cognitive complexity 1 (> max enabled 0)/
	if ok := f(); ok {
		return "it's okay"
	} else {
		return "it's NOT okay!"
	}
}

func h(a, b, c, d, e, f bool) bool { // MATCH /function h has cognitive complexity 2 (> max enabled 0)/
	return a && b && c || d || e && f //FIXME: complexity should be 3
}

func i(a, b, c, d, e, f bool) bool { // MATCH /function i has cognitive complexity 2 (> max enabled 0)/
	result := a && b && c || d || e
	return result
}

func j(a, b, c, d, e, f bool) bool { // MATCH /function j has cognitive complexity 2 (> max enabled 0)/
	result := z(a && !(b && c))
	return result
}

func k(a, b, c, d, e, f bool) bool { // MATCH /function k has cognitive complexity 1 (> max enabled 0)/
	switch expr {
	case cond1:
	case cond2:
	default:
	}

	return result
}

func l() { // MATCH /function l has cognitive complexity 6 (> max enabled 0)/
	for i := 1; i <= max; i++ {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				continue
			}
		}

		total += i
	}
	return total
}

func m() { // MATCH /function m has cognitive complexity 6 (> max enabled 0)/
	if i <= max {
		if j < i {
			if i%j == 0 {
				return 0
			}
		}

		total += i
	}
	return total
}

func n() { // MATCH /function n has cognitive complexity 6 (> max enabled 0)/
	if i > max {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				continue
			}
		}

		total += i
	}
	return total
}

func o() { // MATCH /function o has cognitive complexity 12 (> max enabled 0)/
	if i > max {
		if j < i {
			if i%j == 0 {
				return
			}
		}

		total += i
	}

	if i > max {
		if j < i {
			if i%j == 0 {
				return
			}
		}

		total += i
	}
}

func p() { // MATCH /function p has cognitive complexity 1 (> max enabled 0)/
	switch n := n.(type) {
	case *ast.IfStmt:
		targets := []ast.Node{n.Cond, n.Body, n.Else}
		v.walk(targets...)
		return nil
	case *ast.ForStmt:
		v.walk(n.Body)
		return nil
	case *ast.TypeSwitchStmt:
		v.walk(n.Body)
		return nil
	case *ast.BinaryExpr:
		v.complexity += v.binExpComplexity(n)
		return nil // skip visiting binexp sub-tree (already visited by binExpComplexity)
	}
}

func q() { // MATCH /function q has cognitive complexity 1 (> max enabled 0)/
	for _, t := range targets {
		ast.Walk(v, t)
	}
}

func r() { // MATCH /function r has cognitive complexity 1 (> max enabled 0)/
	select {
	case c <- x:
		x, y = y, x+y
	case <-quit:
		fmt.Println("quit")
		return
	}
}

func s() { // MATCH /function s has cognitive complexity 3 (> max enabled 0)/
FirstLoop:
	for i := 0; i < 10; i++ {
		break
	}
	for i := 0; i < 10; i++ {
		break FirstLoop
	}
}

func t() { // MATCH /function t has cognitive complexity 2 (> max enabled 0)/
FirstLoop:
	for i := 0; i < 10; i++ {
		goto FirstLoop
	}
}

func u() { // MATCH /function u has cognitive complexity 3 (> max enabled 0)/
FirstLoop:
	for i := 0; i < 10; i++ {
		continue
	}
	for i := 0; i < 10; i++ {
		continue FirstLoop
	}
}

func v() { // MATCH /function v has cognitive complexity 2 (> max enabled 0)/
	myFunc := func(b bool) {
		if b {
			// do something
		}
	}
}

func w() { // MATCH /function w has cognitive complexity 3 (> max enabled 0)/
	defer func(b bool) {
		if b {
			// do something
		}
	}(false || true)
}

func sumOfPrimes(max int) int { // MATCH /function sumOfPrimes has cognitive complexity 7 (> max enabled 0)/
	total := 0
OUT:
	for i := 1; i <= max; i++ {
		for j := 2; j < i; j++ {
			if i%j == 0 {
				continue OUT
			}
		}

		total += i
	}
	return total
}
