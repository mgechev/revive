package pkg

import "time"

// these notations are supported by Go, the rule reports them anyway
var (
	// negative values can be used
	// Go will compute 2022-10-28 20:55:54.999999994 +0000 UTC
	_ = time.Date(
		2023,
		-1, // MATCH /time.Date month argument is supposed to be between 1 and 12, found: -1/
		-2, // MATCH /time.Date day argument is supposed to be between 1 and 31, found: -2/
		-3, // MATCH /time.Date hour argument is supposed to be between 0 and 23, found: -3/
		-4, // MATCH /time.Date minute argument is supposed to be between 0 and 59, found: -4/
		-5, // MATCH /time.Date second argument is supposed to be between 0 and 59, found: -5/
		-6, // MATCH /time.Date nanosecond argument is supposed to be between 0 and 999999999, found: -6/
		time.UTC)

	_ = time.Date(
		2023,
		0, // MATCH /time.Date month argument is supposed to be between 1 and 12, found: 0/
		0, // MATCH /time.Date day argument is supposed to be between 1 and 31, found: 0/
		0,
		0,
		0,
		0,
		time.UTC)

	// value out of the expected scale for a field, are OK.
	// Go will compute 2024-02-02 04:11:10 +0000 UTC
	_ = time.Date(
		2023,
		13, // MATCH /time.Date month argument is supposed to be between 1 and 12, found: 13/
		32, // MATCH /time.Date day argument is supposed to be between 1 and 31, found: 32/
		27, // MATCH /time.Date hour argument is supposed to be between 0 and 23, found: 27/
		70, // MATCH /time.Date minute argument is supposed to be between 0 and 59, found: 70/
		70, // MATCH /time.Date second argument is supposed to be between 0 and 59, found: 70/
		0, time.UTC)

	// using gigantic nanoseconds to move into the future
	// Go will compute 2054-09-10 04:50:45 +0000 UTC
	_ = time.Date(2023, 1, 2, 3, 4, 5, 1000000000000000000, time.UTC) // MATCH /time.Date nanosecond argument is supposed to be between 0 and 999999999, found: 1000000000000000000/
)

// Common user errors with dates
// it's out of the scope of the rule
var (
	// there is no February 29th in 2023, Go will compute 2023-03-01 03:04:05 +0000 UTC
	_ = time.Date(2023, 2, 29, 3, 4, 5, 0, time.UTC)

	// June has only 30 days, Go will compute 2023-07-01 03:04:05 +0000 UTC
	_ = time.Date(2023, 6, 31, 3, 4, 5, 0, time.UTC)
)
