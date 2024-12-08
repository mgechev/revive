// Test that we don't ask for comments on sort.Interface methods.

// Package pkg ...
package pkg

// T is ...
type T []int

// Len by itself should get documented.

func (t T) Len() int { return len(t) } // MATCH /exported method T.Len should have comment or be unexported/

// U is ...
type U []int

func (u U) Len() int           { return len(u) }
func (u U) Less(i, j int) bool { return u[i] < u[j] }
func (u U) Swap(i, j int)      { u[i], u[j] = u[j], u[i] }

func (u U) Other() {} // MATCH /exported method U.Other should have comment or be unexported/

// V is ...
type V []int

func (v V) Len() (result int)               { return len(w) }
func (v V) Less(i int, j int) (result bool) { return w[i] < w[j] }
func (v V) Swap(i int, j int)               { v[i], v[j] = v[j], v[i] }

// W is ...
type W []int

func (w W) Swap(i int, j int) {} // MATCH /exported method W.Swap should have comment or be unexported/

// Vv is ...
type Vv []int

func (vv Vv) Len() (result int)               { return len(w) }      // MATCH /exported method Vv.Len should have comment or be unexported/
func (vv Vv) Less(i int, j int) (result bool) { return w[i] < w[j] } // MATCH /exported method Vv.Less should have comment or be unexported/

// X is ...
type X []int

func (x X) Less(i *pkg.tip) (result bool) { return len(x) } // MATCH /exported method X.Less should have comment or be unexported/
