// Package golint comment
package golint

type (
	// O is a shortcut (alias) for map[string]interface{}, e.g. a JSON object.
	O = map[string]interface{}

	// A is shortcut for []O.
	A = []O

	// This Person type is simple
	Person = map[string]interface{}
)

type Foo struct{} // MATCH /exported type Foo should have comment or be unexported/
