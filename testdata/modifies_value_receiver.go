package fixtures

type data struct {
	num   int
	key   *string
	items map[string]bool
}

func (this data) vmethod() {
	this.num = 8 // MATCH /suspicious assignment to a by-value method receiver/
	*this.key = "v.key"
	this.items = make(map[string]bool) // MATCH /suspicious assignment to a by-value method receiver/
	this.items["vmethod"] = true
}

func (a A) Foo() *A {
	a.whatever = true
	return &a
}

func (a A) Clone() (*A, error) {
	a.whatever = true
	return &a, nil
}

// WithBin will set the specific bin path to the builder.
func (b JailerCommandBuilder) WithBin(bin string) JailerCommandBuilder {
	b.bin = bin
	return b
}

func (this data) incrementDecrement() {
	this.num++ // MATCH /suspicious assignment to a by-value method receiver/
	this.num-- // MATCH /suspicious assignment to a by-value method receiver/
	other++
}

func (this data) assignmentOperators() {
	this.num += 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num -= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num *= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num /= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num %= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num &= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num ^= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num |= 1  // MATCH /suspicious assignment to a by-value method receiver/
	this.num >>= 1 // MATCH /suspicious assignment to a by-value method receiver/
	this.num <<= 1 // MATCH /suspicious assignment to a by-value method receiver/
}

type rfoo struct {
	CntP *int
	Cnt  int
	rbar
	RBar  rbar
	PRBar *rbar
}

type rbar struct {
	BCntP *int
	BCnt  int
}

func (foo rfoo) increment() {
	// these are detected
	foo.Cnt++  // MATCH /suspicious assignment to a by-value method receiver/
	foo.BCnt++ // MATCH /suspicious assignment to a by-value method receiver/

	// this one is only a another notation for foo.BCnt++
	foo.rbar.BCnt++ // MATCH /this one should be detected, no?/
	// this on
	foo.RBar.BCnt++ // MATCH /this one should be detected, no?/

	// here, we are updating the pointer of a non-pointer receiver
	// it will lead to nothing, it should be detected, no ?
	*foo.CntP++ // MATCH /what do we want for this one?/

	// same here, it should be detected, no ?
	*foo.rbar.BCntP++ // MATCH /what do we want for this one?/
	*foo.BCntP++      // MATCH /what do we want for this one?/

	// these rely on pointers, they should not be detected, right?
	*foo.RBar.BCntP++  // MATCH /what do we want for this one?/
	*foo.PRBar.BCntP++ // MATCH /what do we want for this one?/
	foo.PRBar.BCnt++   // MATCH /what do we want for this one?/
}
