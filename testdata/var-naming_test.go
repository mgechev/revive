// Test for name linting.

// Package pkg_with_underscores ...
package varnaming_test

var var_name int // MATCH /don't use underscores in Go names; var var_name should be varName/

func Test_ATest()           {}
func Example_AnExample()    {}
func Benchmark_ABenchmark() {}

func Fuzz_AFuzz() {}
