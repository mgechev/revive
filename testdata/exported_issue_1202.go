package golint

import (
	_ "embed"
)

var A []byte // MATCH /exported var A should have comment or be unexported/

var B string // MATCH /exported var B should have comment or be unexported/

//go:embed foo.txt
var C []byte // MATCH /exported var C should have comment or be unexported/

//go:generate pwd
var D string // MATCH /exported var D should have comment or be unexported/

func E() string { // MATCH /exported function E should have comment or be unexported/
	return "E"
}

//nolint:gochecknoglobals
func F() string { // MATCH /exported function F should have comment or be unexported/
	return "F"
}

//nolint:gochecknoglobals
const G = "G" // MATCH /exported const G should have comment or be unexported/
