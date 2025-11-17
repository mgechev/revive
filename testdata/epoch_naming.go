package testdata

import "time"

var (
	creation        = time.Now().Unix() // MATCH /var creation should have a suffix Sec, Second, Seconds/
	creationSeconds = time.Now().Unix()
	createdAtSec    = time.Now().Unix()
	loginTimeMilli  = time.Now().UnixMilli()
	m               = time.Now().UnixMilli() // MATCH /var m should have a suffix Milli, Ms/
	t               = time.Now().UnixNano()  // MATCH /var t should have a suffix Nano, Ns/
	tNano           = time.Now().UnixNano()
	epochNano       = time.Now().UnixNano()

	// Very short but valid names
	sec         = time.Now().Unix()
	ns          = time.Now().UnixNano()
	ms          = time.Now().UnixMilli()
	us          = time.Now().UnixMicro()
	timestampNs = time.Now().UnixNano()
	createdMs   = time.Now().UnixMilli()
	timeSecond  = time.Now().Unix()

	// Edge case 1: UnixMicro() support (Go 1.17+)
	timestampMicro        = time.Now().UnixMicro()
	timestampUs           = time.Now().UnixMicro()
	timestampMicrosecond  = time.Now().UnixMicro()
	timestampMicroseconds = time.Now().UnixMicro()
	badMicroValue         = time.Now().UnixMicro() // MATCH /var badMicroValue should have a suffix Micro, Microsecond, Microseconds, Us/
	invalidUsValue        = time.Now().UnixMicro() // MATCH /var invalidUsValue should have a suffix Micro, Microsecond, Microseconds, Us/

	// Wrong unit suffix - using suffix for different time unit
	badMs         = time.Now().Unix()      // MATCH /var badMs should have a suffix Sec, Second, Seconds/
	wrongSec      = time.Now().UnixMilli() // MATCH /var wrongSec should have a suffix Milli, Ms/
	wrongNano     = time.Now().Unix()      // MATCH /var wrongNano should have a suffix Sec, Second, Seconds/
	confusedMicro = time.Now().UnixNano()  // MATCH /var confusedMicro should have a suffix Nano, Ns/

	// Edge case 3: Case variations (mixed case)
	timeSEC     = time.Now().Unix()      // OK - SEC is valid suffix
	timeNANO    = time.Now().UnixNano()  // OK - NANO is valid suffix
	timeMILLI   = time.Now().UnixMilli() // OK - MILLI is valid suffix
	timeMICRO   = time.Now().UnixMicro() // OK - MICRO is valid suffix
	timeSec     = time.Now().Unix()      // OK - Sec is valid suffix
	timeSeconds = time.Now().Unix()      // OK - Seconds is valid suffix

	// Edge case 6a: Explicit type declarations
	badExplicit     int64 = time.Now().Unix() // MATCH /var badExplicit should have a suffix Sec, Second, Seconds/
	goodSecExplicit int64 = time.Now().Unix() // MATCH /var goodSecExplicit should have a suffix Sec, Second, Seconds/
	explicitSec     int64 = time.Now().Unix()
	explicitSeconds int64 = time.Now().Unix()
	// Edge case 6b: Suffix in middle vs end (suffix must be at the end)
	secInMiddle  int64 = time.Now().Unix()     // MATCH /var secInMiddle should have a suffix Sec, Second, Seconds/
	nanoInMiddle int64 = time.Now().UnixNano() // MATCH /var nanoInMiddle should have a suffix Nano, Ns/
)

func foo() {
	anotherInvalid := time.Now().Unix() // MATCH /var anotherInvalid should have a suffix Sec, Second, Seconds/
	anotherValidSec := time.Now().Unix()
	_ = time.Now().Unix()  // assignment to blank identifier - ignored
	_ := time.Now().Unix() // short declaration with blank identifier - should be ignored
	println(anotherInvalid, anotherValidSec)

	// Edge cases that should NOT be flagged (no variable declaration)
	bar(time.Now().Unix())              // function argument - OK
	_ = []int64{time.Now().UnixMilli()} // slice literal - OK
	baz()

	// Edge case 2: Multiple variables in one statement
	invalidTime := time.Now().Unix()                                // MATCH /var invalidTime should have a suffix Sec, Second, Seconds/
	invalidMilliTime := time.Now().UnixMilli()                      // MATCH /var invalidMilliTime should have a suffix Milli, Ms/
	goodSec, goodMs := time.Now().Unix(), time.Now().UnixMilli()    // Both should pass
	badNanoValue := time.Now().UnixNano()                           // MATCH /var badNanoValue should have a suffix Nano, Ns/
	badMicroValue2 := time.Now().UnixMicro()                        // MATCH /var badMicroValue2 should have a suffix Micro, Microsecond, Microseconds, Us/
	goodNs, goodUs := time.Now().UnixNano(), time.Now().UnixMicro() // Both should pass
	println(invalidTime, invalidMilliTime, goodSec, goodMs, badNanoValue, badMicroValue2, goodNs, goodUs)
}

func bar(input int64) {}

func baz() int64 {
	return time.Now().UnixNano() // return statement - OK
}

// Struct fields should be checked
type Config struct {
	timestamp int64
}

func structTest() {
	c := Config{timestamp: time.Now().Unix()} // struct field - OK (no variable for the timestamp value itself)
}

// Regular assignment (=) should be checked
func regularAssignment() {
	var x int64           // declaration without initialization - OK
	x = time.Now().Unix() // MATCH /var x should have a suffix Sec, Second, Seconds/
	println(x)
}

// Edge case 5: Compound assignments
func compoundAssignment() {
	var counter int64
	counter += time.Now().Unix() // Compound assignment - currently not checked but documented
	println(counter)
}
