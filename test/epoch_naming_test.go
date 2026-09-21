package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestEpochNaming(t *testing.T) {
	testRule(t, "epoch_naming", &rule.EpochNamingRule{})
}
