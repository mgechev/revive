package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestTimeDate(t *testing.T) {
	testRule(t, "time_date_decimal_literal", &rule.TimeDateRule{})
}
