package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestModifiesValRec(t *testing.T) {
	testRule(t, "modifies-value-receiver", &rule.ModifiesValRecRule{})
}
