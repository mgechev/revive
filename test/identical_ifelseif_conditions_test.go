package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestIdenticalIfElseIfConditions(t *testing.T) {
	testRule(t, "identical_ifelseif_conditions", &rule.IdenticalIfElseIfConditionsRule{})
}
