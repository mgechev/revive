package fixtures

import (
	"k8s.io/api/core/v1" // package name is v1
)

func testVer() {
	v1 := "" // Do not warn on this rare case.
	core := "" // MATCH /The name 'core' shadows an import name/
}
