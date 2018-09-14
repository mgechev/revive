package fixtures

func (this data) vmethod() {
	nil := true // MATCH /assignment creates a shadow of built-in identifier nil/
	iota = 1    // MATCH /assignment modifies built-in identifier iota/
}

func append(i, j int) { // MATCH /redefinition of the built-in function append/

}

type Type int16 // MATCH /redefinition of the built-in type Type/

func delete(set []int64, i int) (y []int64) { // MATCH /redefinition of the built-in function delete/
	for j, v := range set {
		if j != i {
			y = append(y, v)
		}
	}
	return
}
