package fixtures

type Foo struct {
	Bar struct { // MATCH /no nested structs are allowed/
	}
}

type Bar struct {
	Baz Baz
}

type Baz struct {
}
