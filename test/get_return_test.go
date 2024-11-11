package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestGetReturn(t *testing.T) {
	testRule(t, "get_return", &rule.GetReturnRule{})
}
