package test

import (
	"testing"

	"github.com/deepsourcelabs/revive/rule"
)

func TestUnusedReceiver(t *testing.T) {
	testRule(t, "unused-receiver", &rule.UnusedReceiverRule{})
}
