package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestModifiesParam(t *testing.T) {
	testRule(t, "modifies_param", &rule.ModifiesParamRule{})
}
