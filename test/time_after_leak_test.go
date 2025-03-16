package test

import (
	"testing"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func TestTimeAfterLeak(t *testing.T) {
	testRule(t, "time_after_leak", &rule.TimeAfterLeak{}, &lint.RuleConfig{})
}

func TestTimeAfterLeakAfterGo1_23(t *testing.T) {
	testRule(t, "go1.23/time_after_leak", &rule.TimeAfterLeak{}, &lint.RuleConfig{})
}
