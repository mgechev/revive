package fixtures

func barSkip() (DummyWriter, error) { // MATCH /barSkip returns interface type: fixtures.DummyWriter/
	return nil, nil
}

func skipDummy() DummyResults {
	return nil
}

func skipDummyFirst() (DummyResults, error) { // MATCH /skipDummyFirst returns interface type: error/
	return nil, nil
}

type Dummy struct{}

func (d Dummy) returnSkippedInterfaceType() DummyResults {
	return nil
}

type DummyResults interface{}
type DummyWriter interface{}
