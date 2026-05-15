package fixtures

func fooAny() any {
	return nil
}

func foo() error {
	return nil
}

func barEmptyInterface() interface{} {
	return nil
}

func bar() DummyReader { // MATCH /bar returns interface type fixtures.DummyReader/
	return nil
}

func fix() Reader { // MATCH /fix returns interface type fixtures.Reader/
	return nil
}

func smile() interface{} {
	return S{}
}

type S struct{}

func (S) Do() {}

type A struct{}

func (a A) B() DummyReader { // MATCH /fixtures.A.B returns interface type fixtures.DummyReader/
	return nil
}

func firstSkip() (DummyWriter, error) { // MATCH /firstSkip returns interface type fixtures.DummyWriter/
	return nil, nil
}

func skipDummyResults() DummyResults {
	return nil
}

func spotMiddle() (any, DummyWriter, error) { // MATCH /spotMiddle returns interface type fixtures.DummyWriter/
	return nil, nil, nil
}

func spotLast() (any, DummyWriter) { // MATCH /spotLast returns interface type fixtures.DummyWriter/
	return nil, nil
}

func skipAll() (DummyResults, any, interface{}, error) {
	return nil, nil, nil, nil
}

type Reader interface {
	Read([]byte) (int, error)
}
type DummyReader interface{}
type DummyResults interface{}
type DummyWriter interface{}

type Dummy struct{}

func (d Dummy) returnSkippedInterfaceType() DummyResults {
	return nil
}
