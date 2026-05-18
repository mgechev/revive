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

func bar() DummyReader { 
	return nil
}

func fix() Reader { 
	return nil
}

func smile() interface{} {
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


