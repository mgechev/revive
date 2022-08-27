// Test of empty-lines.

package fixtures

import "time"

func f1(x *int) bool { // MATCH /extra empty line at the start of a block/

	return x > 2
}

func f2(x *int) bool { // MATCH /extra empty line at the end of a block/
	return x > 2

}

func f3(x *int) bool { // MATCH /extra empty line at the start of a block/

	return x > 2 // MATCH:17 /extra empty line at the end of a block/

}

func f4(x *int) bool {
	// This is fine.
	return x > 2
}

func f5(x *int) bool { // MATCH /extra empty line at the start of a block/

	// This is _not_ fine.
	return x > 2
}

func f6(x *int) bool {
	return x > 2
	// This is fine.
}

func f7(x *int) bool { // MATCH /extra empty line at the end of a block/
	return x > 2
	// This is _not_ fine.

}

func f8(*int) bool {
	if x > 2 { // MATCH /extra empty line at the start of a block/

		return true
	}

	return false
}

func f9(*int) bool {
	if x > 2 { // MATCH /extra empty line at the end of a block/
		return true

	}

	return false
}

func f10(*int) bool { // MATCH /extra empty line at the start of a block/

	if x > 2 {
		return true
	}

	return false
}

func f11(x *int) bool {
	if x > 2 {
		return true
	}
}

func f12(x *int) bool { // MATCH /extra empty line at the end of a block/
	if x > 2 {
		return true
	}

}

func f13(x *int) bool {
	switch {
	case x == 2:
		return false
	}
}

func f14(x *int) bool { // MATCH /extra empty line at the end of a block/
	switch {
	case x == 2:
		return false
	}

}

func f15(x *int) bool {
	switch { // MATCH /extra empty line at the end of a block/
	case x == 2:
		return false

	}
}

func f16(x *int) bool {
	return Query(
		qm("x = ?", x),
	).Execute()
}

func f17(x *int) bool { // MATCH /extra empty line at the end of a block/
	return Query(
		qm("x = ?", x),
	).Execute()

}

func f18(x *int) bool {
	if true {
		if true {
			return true
		}

		// TODO: should we handle the error here?
	}

	return false
}

func w() {
	select {
	case <-time.After(dur):
		// TODO: Handle Ctrl-C is pressed in `mysql` client.
		// return 1 when SLEEP() is KILLed
	}
	return 0, false, nil
}

func x() {
	if tagArray[2] == "req" {
		bit := len(u.reqFields)
		u.reqFields = append(u.reqFields, name)
		reqMask = uint64(1) << uint(bit)
		// TODO: if we have more than 64 required fields, we end up
		// not verifying that all required fields are present.
		// Fix this, perhaps using a count of required fields?
	}

	if err == nil { // No need to refresh if the stream is over or failed.
		// Consider any buffered body data (read from the conn but not
		// consumed by the client) when computing flow control for this
		// stream.
		v := int(cs.inflow.available()) + cs.bufPipe.Len()
		if v < transportDefaultStreamFlow-transportDefaultStreamMinRefresh {
			streamAdd = int32(transportDefaultStreamFlow - v)
			cs.inflow.add(streamAdd)
		}
	}
}

func ShouldNotWarn() {
	// comment

	println()

	// comment
}

// NR test for issue #739
func NotWarnInSingleLineFunction() { println("foo") }
