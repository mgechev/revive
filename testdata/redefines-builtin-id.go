package fixtures

func (this data) vmethod() {
	nil := true // MATCH /assignment creates a shadow of built-in identifier nil/
	iota = 1    // MATCH /assignment modifies built-in identifier iota/
}

func append(i, j int) { // MATCH /redefinition of the built-in function append/

}

type string int16 // MATCH /redefinition of the built-in type string/

func delete(set []int64, i int) (y []int64) { // MATCH /redefinition of the built-in function delete/
	for j, v := range set {
		if j != i {
			y = append(y, v)
		}
	}
	return
}

type any int // MATCH /redefinition of the built-in type any/

func any() {} // MATCH /redefinition of the built-in type any/

var any int // MATCH /redefinition of the built-in type any/

const any = 1 // MATCH /redefinition of the built-in type any/

var i, copy int // MATCH /redefinition of the built-in function copy/

// issue #792
type ()
