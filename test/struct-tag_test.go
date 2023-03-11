package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

// TestStructTag tests struct-tag rule
func TestStructTag(t *testing.T) {
	testRule(t, "struct-tag", &rule.StructTagRule{})
}

func TestStructTagWithUserOptions(t *testing.T) {
	testRule(t, "struct-tag-useroptions", &rule.StructTagRule{}, &lint.RuleConfig{
		Arguments: []interface{}{"json,inline,outline", "bson,gnu"},
	})
}
