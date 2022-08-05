package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestDatarace(t *testing.T) {
	testRule(t, "datarace", &rule.DataRaceRule{})
}
