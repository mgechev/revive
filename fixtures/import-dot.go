// Test that dot imports are flagged.

// Package pkg ...
package pkg

import . "fmt" // MATCH /should not use dot imports/

var _ Stringer // from "fmt"
