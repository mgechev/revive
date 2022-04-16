package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestDatarace(t *testing.T) {
	testRule(t, "datarace", &rule.DataRaceRule{})
}
