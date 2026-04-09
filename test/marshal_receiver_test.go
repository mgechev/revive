package test_test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestMarshalReceiver(t *testing.T) {
	testRule(t, "marshal_receiver", &rule.MarshalReceiverRule{})
}
