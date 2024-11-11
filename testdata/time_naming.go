package testdata

import "time"

var (
	Hour     = time.Hour      // MATCH /var Hour is of type time.Duration; don't use unit-specific suffix "Hour"/
	oneHour  = time.Hour      // MATCH /var oneHour is of type time.Duration; don't use unit-specific suffix "Hour"/
	twoHours = 2 * time.Hour  // MATCH /var twoHours is of type time.Duration; don't use unit-specific suffix "Hours"/
	TenHours = 10 * time.Hour // MATCH /var TenHours is of type time.Duration; don't use unit-specific suffix "Hours"/
	SixHours = 6 * time.Hour  // MATCH /var SixHours is of type time.Duration; don't use unit-specific suffix "Hours"/

	oneMin     = time.Minute     // MATCH /var oneMin is of type time.Duration; don't use unit-specific suffix "Min"/
	Min        = time.Minute     // MATCH /var Min is of type time.Duration; don't use unit-specific suffix "Min"/
	twoMin     = 2 * time.Minute // MATCH /var twoMin is of type time.Duration; don't use unit-specific suffix "Min"/
	SixMin     = 6 * time.Minute // MATCH /var SixMin is of type time.Duration; don't use unit-specific suffix "Min"/
	SixMins    = 6 * time.Minute // MATCH /var SixMins is of type time.Duration; don't use unit-specific suffix "Mins"/
	SixMinutes = 6 * time.Minute // MATCH /var SixMinutes is of type time.Duration; don't use unit-specific suffix "Minutes"/

	oneSec     = time.Second     // MATCH /var oneSec is of type time.Duration; don't use unit-specific suffix "Sec"/
	Sec        = time.Second     // MATCH /var Sec is of type time.Duration; don't use unit-specific suffix "Sec"/
	SixSec     = 6 * time.Second // MATCH /var SixSec is of type time.Duration; don't use unit-specific suffix "Sec"/
	twoSecs    = 2 * time.Second // MATCH /var twoSecs is of type time.Duration; don't use unit-specific suffix "Secs"/
	SixSeconds = 6 * time.Second // MATCH /var SixSeconds is of type time.Duration; don't use unit-specific suffix "Seconds"/
	oneSecond  = time.Second     // MATCH /var oneSecond is of type time.Duration; don't use unit-specific suffix "Second"/
	Second     = time.Second     // MATCH /var Second is of type time.Duration; don't use unit-specific suffix "Second"/

	SixMsec         = 6 * time.Millisecond // MATCH /var SixMsec is of type time.Duration; don't use unit-specific suffix "Msec"/
	oneMsec         = time.Millisecond     // MATCH /var oneMsec is of type time.Duration; don't use unit-specific suffix "Msec"/
	SixMsecs        = 6 * time.Millisecond // MATCH /var SixMsecs is of type time.Duration; don't use unit-specific suffix "Msecs"/
	oneMilli        = time.Millisecond     // MATCH /var oneMilli is of type time.Duration; don't use unit-specific suffix "Milli"/
	SixMillis       = 6 * time.Millisecond // MATCH /var SixMillis is of type time.Duration; don't use unit-specific suffix "Millis"/
	SixMilliseconds = 6 * time.Millisecond // MATCH /var SixMilliseconds is of type time.Duration; don't use unit-specific suffix "Milliseconds"/
	oneMillisecond  = time.Millisecond     // MATCH /var oneMillisecond is of type time.Duration; don't use unit-specific suffix "Millisecond"/
	Millisecond     = time.Millisecond     // MATCH /var Millisecond is of type time.Duration; don't use unit-specific suffix "Millisecond"/

	oneUsec         = 1 * time.Microsecond // MATCH /var oneUsec is of type time.Duration; don't use unit-specific suffix "Usec"/
	twoUsec         = 2 * time.Microsecond // MATCH /var twoUsec is of type time.Duration; don't use unit-specific suffix "Usec"/
	SixUsec         = 6 * time.Microsecond // MATCH /var SixUsec is of type time.Duration; don't use unit-specific suffix "Usec"/
	SixUsecs        = 6 * time.Microsecond // MATCH /var SixUsecs is of type time.Duration; don't use unit-specific suffix "Usecs"/
	twoMicroseconds = 2 * time.Microsecond // MATCH /var twoMicroseconds is of type time.Duration; don't use unit-specific suffix "Microseconds"/
	SixMicroseconds = 6 * time.Microsecond // MATCH /var SixMicroseconds is of type time.Duration; don't use unit-specific suffix "Microseconds"/
	oneMicrosecond  = 1 * time.Microsecond // MATCH /var oneMicrosecond is of type time.Duration; don't use unit-specific suffix "Microsecond"/
	SixMS           = 6 * time.Microsecond // MATCH /var SixMS is of type time.Duration; don't use unit-specific suffix "MS"/
	oneMS           = 1 * time.Microsecond // MATCH /var oneMS is of type time.Duration; don't use unit-specific suffix "MS"/

)
