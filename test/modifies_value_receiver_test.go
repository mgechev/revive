package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestModifiesValRec(t *testing.T) {
	testRule(t, "modifies_value_receiver", &rule.ModifiesValRecRule{})
}

// TestModifiesValRec_typeCheckDoesNotPanic checks that go/types does not panic on invalid
// nested indexed composite literals (issue #59), since Package.TypeCheck does not recover.
func TestModifiesValRec_typeCheckDoesNotPanic(t *testing.T) {
	testRule(t, "modifies_value_receiver_issue_59", &rule.ModifiesValRecRule{})
}
