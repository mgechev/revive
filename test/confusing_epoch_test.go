package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestConfusingEpoch(t *testing.T) {
	testRule(t, "confusing_epoch", &rule.ConfusingEpochRule{})
}
