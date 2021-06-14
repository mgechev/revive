package fixtures

type Foo struct {
	Bar struct { // MATCH /no nested structs are allowed/
		Baz struct { // MATCH /no nested structs are allowed/
			b   bool
			Qux struct { // MATCH /no nested structs are allowed/
				b bool
			}
		}
	}
}

type Quux struct {
	Quuz Quuz
}

type Quuz struct {
}

func waldo() (s struct{ b bool }) { return s }

func fred() interface{} {
	s := struct {
		b bool
		t struct { // MATCH /no nested structs are allowed/
			b bool
		}
	}{}

	return s
}
