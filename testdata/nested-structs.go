package fixtures // MATCH /no nested structs are allowed, got 1/

type Foo struct {
	Bar struct {
	}
}

type Bar struct {
	Baz Baz
}

type Baz struct {
}
