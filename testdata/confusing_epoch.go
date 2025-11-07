package testdata

import "time"

var (
	creation        = time.Now().Unix() // MATCH /variable 'creation' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	creationSeconds = time.Now().Unix()
	createdAtSec    = time.Now().Unix()
	loginTimeMilli  = time.Now().UnixMilli()
	m               = time.Now().UnixMilli() // MATCH /variable 'm' initialized with UnixMilli() should have a name containing one of [Milli Ms]/
	t               = time.Now().UnixNano()  // MATCH /variable 't' initialized with UnixNano() should have a name containing one of [Nano Ns]/
	tNano           = time.Now().UnixNano()
	epochNano       = time.Now().UnixNano() // OK - "Nano" suffix at end makes this valid

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
	badMicroValue         = time.Now().UnixMicro() // MATCH /variable 'badMicroValue' initialized with UnixMicro() should have a name containing one of [Micro Microsecond Microseconds Us]/
	invalidUsValue        = time.Now().UnixMicro() // MATCH /variable 'invalidUsValue' initialized with UnixMicro() should have a name containing one of [Micro Microsecond Microseconds Us]/

	// Wrong unit suffix - using suffix for different time unit
	badMs         = time.Now().Unix()      // MATCH /variable 'badMs' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	wrongSec      = time.Now().UnixMilli() // MATCH /variable 'wrongSec' initialized with UnixMilli() should have a name containing one of [Milli Ms]/
	wrongNano     = time.Now().Unix()      // MATCH /variable 'wrongNano' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	confusedMicro = time.Now().UnixNano()  // MATCH /variable 'confusedMicro' initialized with UnixNano() should have a name containing one of [Nano Ns]/

	// Edge case 3: Case variations (mixed case)
	timeSEC     = time.Now().Unix()      // OK - SEC is valid suffix
	timeNANO    = time.Now().UnixNano()  // OK - NANO is valid suffix
	timeMILLI   = time.Now().UnixMilli() // OK - MILLI is valid suffix
	timeMICRO   = time.Now().UnixMicro() // OK - MICRO is valid suffix
	timeSec     = time.Now().Unix()      // OK - Sec is valid suffix
	timeSeconds = time.Now().Unix()      // OK - Seconds is valid suffix

	// Edge case 6: Suffix in middle vs end (suffix must be at the end)
	badExplicit     int64 = time.Now().Unix() // MATCH /variable 'badExplicit' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	goodSecExplicit int64 = time.Now().Unix() // MATCH /variable 'goodSecExplicit' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	explicitSec     int64 = time.Now().Unix()
	explicitSeconds int64 = time.Now().Unix()
	secInMiddle     int64 = time.Now().Unix()     // MATCH /variable 'secInMiddle' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	nanoInMiddle    int64 = time.Now().UnixNano() // MATCH /variable 'nanoInMiddle' initialized with UnixNano() should have a name containing one of [Nano Ns]/
)

func foo() {
	anotherInvalid := time.Now().Unix() // MATCH /variable 'anotherInvalid' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	anotherValidSec := time.Now().Unix()
	_ = time.Now().Unix()  // assignment to blank identifier - ignored
	_ := time.Now().Unix() // short declaration with blank identifier - should be ignored
	println(anotherInvalid, anotherValidSec)

	// Edge cases that should NOT be flagged (no variable declaration)
	bar(time.Now().Unix())              // function argument - OK
	_ = []int64{time.Now().UnixMilli()} // slice literal - OK
	baz()

	// Edge case 2: Multiple variables in one statement
	invalidTime := time.Now().Unix()                                // MATCH /variable 'invalidTime' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	invalidMilliTime := time.Now().UnixMilli()                      // MATCH /variable 'invalidMilliTime' initialized with UnixMilli() should have a name containing one of [Milli Ms]/
	goodSec, goodMs := time.Now().Unix(), time.Now().UnixMilli()    // Both should pass
	badNanoValue := time.Now().UnixNano()                           // MATCH /variable 'badNanoValue' initialized with UnixNano() should have a name containing one of [Nano Ns]/
	badMicroValue := time.Now().UnixMicro()                         // MATCH /variable 'badMicroValue' initialized with UnixMicro() should have a name containing one of [Micro Microsecond Microseconds Us]/
	goodNs, goodUs := time.Now().UnixNano(), time.Now().UnixMicro() // Both should pass
	println(invalidTime, invalidMilliTime, goodSec, goodMs, badNanoValue, badMicroValue, goodNs, goodUs)
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
	x = time.Now().Unix() // MATCH /variable 'x' initialized with Unix() should have a name containing one of [Sec Second Seconds]/
	println(x)
}

// Edge case 5: Compound assignments
func compoundAssignment() {
	var counter int64
	counter += time.Now().Unix() // Compound assignment - currently not checked but documented
	println(counter)
}
