package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestStructTag(t *testing.T) {
	testRule(t, "struct_tag", &rule.StructTagRule{})
}

func TestStructTagWithUserOptions(t *testing.T) {
	testRule(t, "struct_tag_user_options", &rule.StructTagRule{}, &lint.RuleConfig{
		Arguments: []any{"json,inline,outline", "bson,gnu"},
	})
}
