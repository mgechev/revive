// Package golint comment
package golint

// Test case for the configuration option tp replace the word "stutters" by "repetitive" failure messages

//  GolintRepetitive is a dummy function
func GolintRepetitive() {} // MATCH /func name will be used as golint.GolintRepetitive by other packages, and that is repetitive; consider calling this Repetitive/
