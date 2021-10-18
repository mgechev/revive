package fixtures

func conditionalExpr(x, y bool) bool {
	equal := x == y // should not match, not AND or OR operators
	if x || y {     // should not match, no caller
		return true
	}
	or := caller(x, y) || y // MATCH /for better performance 'caller(x, y) || y' might be rewritten as 'y || caller(x, y)'/
	if caller(x, y) || y {  // MATCH /for better performance 'caller(x, y) || y' might be rewritten as 'y || caller(x, y)'/
		return true
	}

	switch {
	case x == y:
		return y
	case caller(x, y) && y: // MATCH /for better performance 'caller(x, y) && y' might be rewritten as 'y && caller(x, y)'/
		return x
	}

	complexExpr := caller(caller(x, y) && y, y) || y
	// MATCH:20 /for better performance 'caller(caller(x, y) && y, y) || y' might be rewritten as 'y || caller(caller(x, y) && y, y)'/
	// MATCH:20 /for better performance 'caller(x, y) && y' might be rewritten as 'y && caller(x, y)'/

	noSwap := caller(x, y) || (x && caller(y, x)) // should not match, caller in the right operand

	callRight := caller(x, y) && (x && caller(y, x)) // should not match, caller in the right operand
	return caller(x, y) && y                         // MATCH /for better performance 'caller(x, y) && y' might be rewritten as 'y && caller(x, y)'/
}

func caller(x, y bool) bool {
	return true
}
