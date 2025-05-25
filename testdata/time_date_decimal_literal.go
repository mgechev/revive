package pkg

import "time"

// these should match and be reported as invalid usage of time.Date
// the rule should suggest to use a decimal number
var (
	// All these examples refers to the same date
	// 2023-01-02 03:04:05.000000006 +0000 UTC

	// this is an invalid usage
	// people are using this format because a date is easier to read
	// but here it leads to use octal numbers0
	_ = time.Date(
		2023,
		01,        // MATCH /use decimal digits for time.Date month argument: octal notation with leading zero found: use 1 instead of 01/
		02,        // MATCH /use decimal digits for time.Date day argument: octal notation with leading zero found: use 2 instead of 02/
		03,        // MATCH /use decimal digits for time.Date hour argument: octal notation with leading zero found: use 3 instead of 03/
		04,        // MATCH /use decimal digits for time.Date minute argument: octal notation with leading zero found: use 4 instead of 04/
		05,        // MATCH /use decimal digits for time.Date second argument: octal notation with leading zero found: use 5 instead of 05/
		000000006, // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with padding zeroes found: use 6 instead of 000000006/
		time.UTC)

	// the following one could have been written by someone who is not aware of the issue
	// Please note, there are multiple issues on the same line
	//
	// use special syntax to match multiple issues being reported on the same line
	// MATCH:34 /use decimal digits for time.Date month argument: octal notation with leading zero found: use 1 instead of 01/
	// MATCH:34 /use decimal digits for time.Date day argument: octal notation with leading zero found: use 2 instead of 02/
	// MATCH:34 /use decimal digits for time.Date hour argument: octal notation with leading zero found: use 3 instead of 03/
	// MATCH:34 /use decimal digits for time.Date minute argument: octal notation with leading zero found: use 4 instead of 04/
	// MATCH:34 /use decimal digits for time.Date second argument: octal notation with leading zero found: use 5 instead of 05/
	// MATCH:34 /use decimal digits for time.Date nanosecond argument: octal notation with padding zeroes found: use 6 instead of 000000006/
	_ = time.Date(2023, 01, 02, 03, 04, 05, 000000006, time.UTC)
)

// gofumpt formats legacy non-decimal notation to new one non-decimal notation
// it transforms 01, 02, 03, 04, 05, 06, and 07 to 0o1, 0o2, 0o3, 0o4, 0o5, 0o6, and 0o7
// but here with time.Date it doesn't make sense. This is the main reason why the rule was created.
var (
	_ = time.Date(2023, 0o1, 2, 3, 4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date month argument: octal notation found: use 1 instead of 0o1/
	_ = time.Date(2023, 1, 0o2, 3, 4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date day argument: octal notation found: use 2 instead of 0o2/
	_ = time.Date(2023, 1, 2, 0o3, 4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date hour argument: octal notation found: use 3 instead of 0o3/
	_ = time.Date(2023, 1, 2, 3, 0o4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date minute argument: octal notation found: use 4 instead of 0o4/
	_ = time.Date(2023, 1, 2, 3, 4, 0o5, 6, time.UTC) // MATCH /use decimal digits for time.Date second argument: octal notation found: use 5 instead of 0o5/
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0o6, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation found: use 6 instead of 0o6/
)

// padding with zeroes can lead to errors
var (
	_ = time.Date(2023, 1, 2, 3, 4, 5, 00, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with leading zero found: use 0 instead of 00/
	_ = time.Date(2023, 1, 2, 3, 4, 5, 01, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with leading zero found: use 1 instead of 01/

	_ = time.Date(2023, 1, 2, 3, 4, 5, 00000000, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with padding zeroes found: use 0 instead of 00000000/
	_ = time.Date(2023, 1, 2, 3, 4, 5, 00000006, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with padding zeroes found: use 6 instead of 00000006/
	_ = time.Date(2023, 1, 2, 3, 4, 5, 00123456, time.UTC) // MATCH /use decimal digits for time.Date nanosecond argument: octal notation with padding zeroes found: choose between 123456 and 42798 (decimal value of 123456 octal value)/
)

// hypothetical examples based on other number notations
// https://go.dev/ref/spec#Integer_literals
// these should match and be reported as invalid usage of time.Date
var (
	_ = time.Date(
		0x7e7,    // MATCH /use decimal digits for time.Date year argument: hexadecimal notation found: use 2023 instead of 0x7e7/
		0b1,      // MATCH /use decimal digits for time.Date month argument: binary notation found: use 1 instead of 0b1/
		0x_2,     // MATCH /use decimal digits for time.Date day argument: hexadecimal notation found: use 2 instead of 0x_2/
		1_3,      // MATCH /use decimal digits for time.Date hour argument: alternative notation found: use 13 instead of 1_3/
		1e1,      // MATCH /use decimal digits for time.Date minute argument: exponential notation found: use 10 instead of 1e1/
		0.,       // MATCH /use decimal digits for time.Date second argument: float literal found: use 0 instead of 0./
		0x1.Fp+6, // MATCH /use decimal digits for time.Date nanosecond argument: float literal found: use 124 instead of 0x1.Fp+6/
		time.UTC)
)

// here we are checking that we also detect non-decimal notation in methods/functions/lambdas
func _() {
	_ = time.Date(2023, 01, 2, 3, 4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date month argument: octal notation with leading zero found: use 1 instead of 01/

	_ = func() time.Time {
		return time.Date(2023, 01, 2, 3, 4, 5, 6, time.UTC) // MATCH /use decimal digits for time.Date month argument: octal notation with leading zero found: use 1 instead of 01/
	}
}

// these should never match
var (
	_ = time.Date(2023, 1, 2, 3, 4, 5, 1234567, time.UTC)
	_ = time.Date(2023, time.January, 2, 3, 4, 5, 1234567, time.UTC)
	_ = time.Date(2023, 10, 10, 10, 10, 10, 100, time.UTC)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 6, time.UTC)
)

// uncommon notations
// these notations are supported by Go, the rule should not report it
// it's out of the scope of the rule
var (
	// negative values can be used
	// Go will compute 2022-10-28 20:55:54.999999994 +0000 UTC
	_ = time.Date(2023, -1, -2, -3, -4, -5, -6, time.UTC)

	// value out of the expected scale for a field, are OK.
	// Go will compute 2024-02-02 04:11:10 +0000 UTC
	_ = time.Date(2023, 13, 32, 27, 70, 70, 0, time.UTC)

	// using gigantic nanoseconds to move into the future
	// Go will compute 2054-09-10 04:50:45 +0000 UTC
	_ = time.Date(2023, 1, 2, 3, 4, 5, 1000000000000000000, time.UTC)
)

// Common user errors with dates
// it's out of the scope of the rule
var (
	// there is no February 29th in 2023, Go will compute 2023-03-01 03:04:05 +0000 UTC
	_ = time.Date(2023, 2, 29, 3, 4, 5, 0, time.UTC)

	// June has only 30 days, Go will compute 2023-07-01 03:04:05 +0000 UTC
	_ = time.Date(2023, 6, 31, 3, 4, 5, 0, time.UTC)
)
