// Test of confusing-naming rule.

// Package pkg ...
package pkg

func aglobal() { // MATCH /Function 'aglobal' differs only by capitalization to other function in the same package/
}
