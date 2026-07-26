package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestMultilineIfInit(t *testing.T) {
	testRule(t, "multiline_if_init", &rule.MultilineIfInitRule{})
}
