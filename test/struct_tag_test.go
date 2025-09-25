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
		Arguments: []any{
			"json,inline,outline",
			"bson,gnu",
			"url,myURLOption",
			"datastore,myDatastoreOption",
			"mapstructure,myMapstructureOption",
			"validate,displayName",
			"toml,unknown",
			"spanner,mySpannerOption",
			"codec,myCodecOption",
			"cbor,myCborOption",
		},
	})
}

func TestStructTagWithOmittedTags(t *testing.T) {
	testRule(t, "struct_tag_user_options_omit", &rule.StructTagRule{}, &lint.RuleConfig{
		Arguments: []any{
			"!validate",
			"!toml",
			"json,inline,outline",
			"bson,gnu",
			"url,myURLOption",
			"datastore,myDatastoreOption",
			"mapstructure,myMapstructureOption",
			"spanner,mySpannerOption",
		},
	})
}

func TestStructTagAfterGo1_24(t *testing.T) {
	testRule(t, "go1.24/struct_tag", &rule.StructTagRule{})
}
