package test

import (
	"testing"

	"github.com/mgechev/revive/rule"
)

func TestUnnecessaryStmt(t *testing.T) {
	testRule(t, "unnecessary_stmt", &rule.UnnecessaryStmtRule{})
}
