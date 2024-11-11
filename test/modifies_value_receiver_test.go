package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestModifiesValRec(t *testing.T) {
	testRule(t, "modifies_value_receiver", &rule.ModifiesValRecRule{})
}
