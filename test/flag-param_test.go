package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestFlagParam(t *testing.T) {
	testRule(t, "flag-param", &rule.FlagParamRule{})
}
