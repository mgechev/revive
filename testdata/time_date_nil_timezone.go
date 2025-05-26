package pkg

import "time"

var (
	// a nil timezone will panic at runtime
	// this is an invalid usage
	// it should be reported as an error
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, nil) // MATCH /time.Date timezone argument cannot be nil, it would panic on runtime/
)

func _() {
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, nil) // MATCH /time.Date timezone argument cannot be nil, it would panic on runtime/
}

func _() {
	// this is a valid usage
	// it should not be reported as an error
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, time.UTC)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, time.Local)

	loc := time.LoadLocation("Europe/Paris")
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, loc)

	loc := time.FixedZone("UTC-8", -8*60*60)
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, loc)

	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, time.FixedZone("UTC-8", -8*60*60))
}

// this would be difficult to detect
// and are for now not reported
// even if they panic at runtime
func _() {
	var loc *time.Location
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, loc)

	loc, _ = time.LoadLocation("whatever")
	_ = time.Date(2023, 1, 2, 3, 4, 5, 0, loc)
}
