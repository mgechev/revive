package fixtures

import (
	"k8s.io/api/core/v1" // package name is v1
)

func testVer() {
	// Do not warn on this rare case.
	v1 := ""
}
