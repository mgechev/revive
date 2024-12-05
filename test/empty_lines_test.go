package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestEmptyLines(t *testing.T) {
	testRule(t, "empty_lines", &rule.EmptyLinesRule{})
}
