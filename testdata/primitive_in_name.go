package testdata

const maxInt = 100 // MATCH /avoid primitive type in name/

const retryCountInt int = 3 // MATCH /avoid primitive type in name/

const MaxCountLimit = 10 // untyped int constant, but no primitive word segment in the name

func PrimitiveInNameExample() {
	a := 5
	b := 10
	sumInt := a + b // MATCH /avoid primitive type in name/
	_ = sumInt

	interfaceLen := 5 // "Int" is a substring of "interface", not a distinct word segment
	_ = interfaceLen

	sprintfCount := 5 // "Int" is not a segment of "Sprintf"
	_ = sprintfCount

	sum_int := a + b // underscore-separated names aren't split into words; var-naming already flags snake_case separately
	_ = sum_int

	var isValidBool bool = true // MATCH /avoid primitive type in name/
	_ = isValidBool

	var nameString string // MATCH /avoid primitive type in name/
	_ = nameString

	var avgFloat32 float32 // MATCH /avoid primitive type in name/
	_ = avgFloat32

	var avgFloat64 float64 // MATCH /avoid primitive type in name/
	_ = avgFloat64

	var firstRune rune // MATCH /avoid primitive type in name/
	_ = firstRune

	var firstByte byte // MATCH /avoid primitive type in name/
	_ = firstByte

	var byteSlice []byte // name contains "Byte" but the variable's type is a slice, not a primitive
	_ = byteSlice

	var c complex128 = 1 + 2i // regression: complex is not in the primitive word list, and must not panic
	_ = c

	var nameStr string // MATCH /avoid primitive type in name/
	_ = nameStr

	var retryNum int // MATCH /avoid primitive type in name/
	_ = retryNum

	var firstChar rune // MATCH /avoid primitive type in name/
	_ = firstChar

	var verboseFlag bool // MATCH /avoid primitive type in name/
	_ = verboseFlag

	var avgFloat float32 // MATCH /avoid primitive type in name/
	_ = avgFloat

	var totalFloat float64 // MATCH /avoid primitive type in name/
	_ = totalFloat
}
