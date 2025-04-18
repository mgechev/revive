package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceRepeatedArgTypeStyleDefault(t *testing.T) {
	testRule(t, "enforce_repeated_arg_type_style_default", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{})
}

func TestEnforceRepeatedArgTypeStyleShort(t *testing.T) {
	testRule(t, "enforce_repeated_arg_type_style_short_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"short"},
	})
	testRule(t, "enforce_repeated_arg_type_style_short_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"short"},
	})

	testRule(t, "enforce_repeated_arg_type_style_short_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `short`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_short_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"func-arg-style": `short`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_short_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `short`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_short_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"func-ret-val-style": `short`,
			},
		},
	})
}

func TestEnforceRepeatedArgTypeStyleFull(t *testing.T) {
	testRule(t, "enforce_repeated_arg_type_style_full_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"full"},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"full"},
	})

	testRule(t, "enforce_repeated_arg_type_style_full_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
}

func TestEnforceRepeatedArgTypeStyleMixed(t *testing.T) {
	testRule(t, "enforce_repeated_arg_type_style_full_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `any`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})

	testRule(t, "enforce_repeated_arg_type_style_full_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `any`,
				"funcRetValStyle": `full`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_full_return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})

	testRule(t, "enforce_repeated_arg_type_style_mixed_full_short", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})
	testRule(t, "enforce_repeated_arg_type_style_mixed_short_full", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})
}
