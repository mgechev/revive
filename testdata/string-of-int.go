// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package fixtures

type A string
type B = string
type C int
type D = uintptr

func StringTest() {
	var (
		i int
		j rune
		k byte
		l C
		m D
		n = []int{0, 1, 2}
		o struct{ x int }
	)
	const p = 0
	_ = string(i) // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
	_ = string(j)
	_ = string(k)
	_ = string(p)    // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
	_ = A(l)         // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
	_ = B(m)         // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
	_ = string(n[1]) // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
	_ = string(o.x)  // MATCH /dubious conversion of an integer into a string, use strconv.Itoa/
}
