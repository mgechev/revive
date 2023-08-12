// should fail if upperCaseConst = false (by default) #851, #865

package fixtures

const SOME_CONST_2 = 1          // MATCH /don't use ALL_CAPS in Go names; use CamelCase/
const _SOME_PRIVATE_CONST_2 = 2 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/

const (
	SOME_CONST_3          = 3 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/
	_SOME_PRIVATE_CONST_3 = 3 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/
)
