package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnexportedNaming(t *testing.T) {
	testRule(t, "unexported_naming", &rule.UnexportedNamingRule{})
}
