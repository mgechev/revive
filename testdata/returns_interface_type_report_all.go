package fixtures

func fooAny() any {// MATCH /fooAny returns interface type any/
	return nil
}

func foo() error {// MATCH /foo returns interface type error/
	return nil
}

func barEmptyInterface() interface{} {// MATCH /barEmptyInterface returns interface type interface{}/
	return nil
}

func bar() DummyReader { // MATCH /bar returns interface type fixtures.DummyReader/
	return nil
}

func fix() Reader { // MATCH /fix returns interface type fixtures.Reader/
	return nil
}

func smile() interface{} {// MATCH /smile returns interface type interface{}/
	return S{}
}

type S struct{}

func (S) Do() {}

type Reader interface {
	Read([]byte) (int, error)
}
type DummyReader interface{}
type DummyResults interface{}
type DummyWriter interface{}

type Dummy struct{}


