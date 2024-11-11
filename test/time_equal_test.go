package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

// TestTimeEqual rule.
func TestTimeEqual(t *testing.T) {
	testRule(t, "time_equal", &rule.TimeEqualRule{})
}
