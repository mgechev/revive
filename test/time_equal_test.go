package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestTimeEqual(t *testing.T) {
	testRule(t, "time_equal", &rule.TimeEqualRule{})
}
