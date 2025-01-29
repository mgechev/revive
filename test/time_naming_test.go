package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestTimeNaming(t *testing.T) {
	testRule(t, "time_naming", &rule.TimeNamingRule{})
}
