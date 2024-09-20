package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestDatarace(t *testing.T) {
	testRule(t, "datarace", &rule.DataRaceRule{})
}

func TestDataraceAfterGo1_22(t *testing.T) {
	testRule(t, "go1.22/datarace", &rule.DataRaceRule{})
}
