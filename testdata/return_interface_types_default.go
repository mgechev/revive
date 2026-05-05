package fixtures

func foo() error {
	return nil
}

func bar() DummyReader { // MATCH /bar returns interface type: fixtures.DummyReader/
	return nil
}

type A struct{}

func (a A) B() DummyReader { // MATCH /fixtures.A.B returns interface type: fixtures.DummyReader/
	return nil
}

type DummyReader interface{}
