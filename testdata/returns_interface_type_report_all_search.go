package fixtures

func fooAnySearch() any {
	return nil
}

func foo() error {// MATCH /foo returns interface type error/
	return nil
}


