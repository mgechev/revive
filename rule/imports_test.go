package rule

import (
	"testing"

	"github.com/mgechev/revive/rule"
	"github.com/mgechev/revive/testutil"
)

func TestImports_Failure(t *testing.T) {
	t.Parallel()

	program := `
	package foo
	
	import (
		"fmt"
		[@f1]. "path"[/@f1]
	)
	     
	`
	testutil.AssertFailures(t, program, &ImportsRule{}, rule.Arguments{})
}

func TestImports(t *testing.T) {
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
	testutil.AssertSuccess(t, program, &ImportsRule{}, rule.Arguments{})
}

func TestImports_SkipTesting(t *testing.T) {
	t.Parallel()

	program := `
	package main
	
	import (
		"fmt"
	
		/* MATCH /blank import/ */ _ "os"
	
		/* MATCH /blank import/ */ _ "net/http"
		. "path"
	)
	     
	`
	testutil.AssertSuccessWithName(t, program, "foo_test.go", &ImportsRule{}, rule.Arguments{})
}
