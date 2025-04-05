package golint

// Deprecated: this is deprecated, use math.PI instead
const PI = 3.14 // MATCH /exported const PI should have comment or be unexported/

// Deprecated: this is deprecated
var Buffer []byte // MATCH /exported var Buffer should have comment or be unexported/

// Deprecated: this is deprecated, use min instead
func Min(a, b int) int { // MATCH /exported function Min should have comment or be unexported/
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum one of two integers.
// Deprecated: this is deprecated, use max instead
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Deprecated: this is deprecated
type Number float64 // MATCH /exported type Number should have comment or be unexported/

// Name is a type that represents a name.
type Name string

// Greet returns a greeting for the name.
func (n Name) Greet() string {
	return "Hello, " + string(n)
}

// Deprecated: this is deprecated, use Name.ToString instead
func (n Name) ToString() string { // MATCH /exported method Name.ToString should have comment or be unexported/
	return string(n)
}

// String returns the string representation of the name.
// Deprecated: this is deprecated, use Name.Greet instead
func (n Name) String() string {
	return string(n)
}
