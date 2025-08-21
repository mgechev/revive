package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestSecretsSerialization(t *testing.T) {
	testRule(t, "secrets_serialization_default", &rule.SecretsSerializationRule{}, &lint.RuleConfig{})
	testRule(t, "secrets_serialization_custom", &rule.SecretsSerializationRule{}, &lint.RuleConfig{
		Arguments: []any{[]any{"email", "SSN"}},
	})
}
