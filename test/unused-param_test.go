package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnusedParam(t *testing.T) {
	testRule(t, "unused-param", &rule.UnusedParamRule{})
}
