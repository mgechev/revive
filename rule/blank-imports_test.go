package rule

import (
	"testing"

	"github.com/mgechev/revive/rule"
	"github.com/mgechev/revive/testutil"
)

func TestBlankImports(t *testing.T) {
	t.Parallel()

	program := `
	package foo
	
	import (
		"fmt"
	
		/* MATCH /blank import/ */ [@f1]_ "os"[/@f1]
	
		/* MATCH /blank import/ */ [@f2]_ "net/http"[/@f2]
		_ "path"
	)
	     
	`
	testutil.AssertFailures(t, program, &BlankImportsRule{}, rule.Arguments{})
}

func TestBlankImports_ShouldSkipMain(t *testing.T) {
	t.Parallel()

	program := `
	package main
	
	import (
		"fmt"
	
		/* MATCH /blank import/ */ _ "os"
	
		/* MATCH /blank import/ */ _ "net/http"
		_ "path"
	)
	     
	`
	testutil.AssertSuccess(t, program, &BlankImportsRule{}, rule.Arguments{})
}
