package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestUnexportednaming(t *testing.T) {
	testRule(t, "unexported-naming", &rule.UnexportedNamingRule{})
}
