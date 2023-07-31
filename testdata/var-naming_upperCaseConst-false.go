// should fail if upperCaseConst = false (by default) #851

package fixtures

const SOME_CONST_2 = 1 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/

const (
	SOME_CONST_3 = 3 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/
	// DUE TO LEGACY - names less 5 symbols without undescores are not treated as UPPER_CASE
	// It's strange, but VERYLARGENAME is not linted too
	VER = 0
)
