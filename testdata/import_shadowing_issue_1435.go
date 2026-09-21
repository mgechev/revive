package fixtures

import (
	"v1"
	"V1"
	"v12345"
	"math/rand/v2"
	randv2 "math/rand/v2"
)

func testVer() {
	v1 := ""     // MATCH /The name 'v1' shadows an import name/
	V1 := ""     // MATCH /The name 'V1' shadows an import name/
	v12345 := "" // MATCH /The name 'v12345' shadows an import name/
	v2 := ""
	V2 := ""
	rand := ""   // MATCH /The name 'rand' shadows an import name/
	randv2 := "" // MATCH /The name 'randv2' shadows an import name/
}
