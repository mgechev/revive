package fixtures

func foo() error {
	return nil
}

func bar() DummyReader { // MATCH /bar returns interface type: fixtures.DummyReader/
	return nil
}

func barEmptyInterface() interface{} {
	return nil
}

func fix() Reader { // MATCH /fix returns interface type: fixtures.Reader/
	return nil
}

func smile() interface{} {
	return S{}
}

type S struct{}

func (S) Do() {}

type A struct{}

func (a A) B() DummyReader { // MATCH /fixtures.A.B returns interface type: fixtures.DummyReader/
	return nil
}

type DummyReader interface{}

type Reader interface {
	Read([]byte) (int, error)
}
