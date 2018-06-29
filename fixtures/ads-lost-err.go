package fixtures

import "errors"

func foo(a, b, c, d int) {
	errors.New("aaa")
	errors.New(errors.InternalError, errors.MessageOption("nice error message "+err.Error())) // MATCH /original error is lost, consider using errors.NewFromError/
	errors.MessageOption("nice error message " + err.Error())
}
