package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestEnforceRepeatedArgTypeStyleShort(t *testing.T) {
	testRule(t, "enforce-repeated-arg-type-style-short-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"short"},
	})
	testRule(t, "enforce-repeated-arg-type-style-short-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"short"},
	})

	testRule(t, "enforce-repeated-arg-type-style-short-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `short`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-short-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `short`,
			},
		},
	})
}

func TestEnforceRepeatedArgTypeStyleFull(t *testing.T) {
	testRule(t, "enforce-repeated-arg-type-style-full-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"full"},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{"full"},
	})

	testRule(t, "enforce-repeated-arg-type-style-full-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
}

func TestEnforceRepeatedArgTypeStyleMixed(t *testing.T) {
	testRule(t, "enforce-repeated-arg-type-style-full-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle": `full`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `any`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-args", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})

	testRule(t, "enforce-repeated-arg-type-style-full-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcRetValStyle": `full`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `any`,
				"funcRetValStyle": `full`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-full-return", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})

	testRule(t, "enforce-repeated-arg-type-style-mixed-full-short", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `full`,
				"funcRetValStyle": `short`,
			},
		},
	})
	testRule(t, "enforce-repeated-arg-type-style-mixed-short-full", &rule.EnforceRepeatedArgTypeStyleRule{}, &lint.RuleConfig{
		Arguments: []any{
			map[string]any{
				"funcArgStyle":    `short`,
				"funcRetValStyle": `full`,
			},
		},
	})
}
