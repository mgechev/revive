package fixtures

import (
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/mgechev/revive/lint"
)

func f0(param int) {
	param := param
}

func f1(param int) { // MATCH /parameter 'param' seems to be unused, consider removing or renaming it as _/
	if param := fn(); predicate(param) {
		// do stuff
	}
}

func f2(param int) { // MATCH /parameter 'param' seems to be unused, consider removing or renaming it as _/
	switch param := fn(); param {
	default:
		
	}
}

func f3(param myStruct) { 
	a:= param.field
}

func f4(param myStruct, c int) { // MATCH /parameter 'c' seems to be unused, consider removing or renaming it as _/ 
	param.field = "aString"
	param.c = "sss" 
}

func f5(a int, _ float) { // MATCH /parameter 'a' seems to be unused, consider removing or renaming it as _/
	fmt.Printf("Hello, Golang\n")
	{
		if true {
			a := 2
			b := a
		}
	}
}

func f6(_ float, c string) { // MATCH /parameter 'c' seems to be unused, consider removing or renaming it as _/
	fmt.Printf("Hello, Golang\n")
	c := 1
}

func assertSuccess(t *testing.T, baseDir string, fi os.FileInfo, src []byte, rules []lint.Rule, config map[string]lint.RuleConfig) error { // MATCH /parameter 'src' seems to be unused, consider removing or renaming it as _/
	l := lint.New(func(file string) ([]byte, error) {
		return ioutil.ReadFile(baseDir + file)
	})

	ps, err := l.Lint([][]string{[]string{fi.Name()}}, rules, lint.Config{
		Rules: config,
	})
	if err != nil {
		return err
	}

	failures := ""
	for p := range ps {
		failures += p.Failure
	}
	if failures != "" {
		t.Errorf("Expected the rule to pass but got the following failures: %s", failures)
	}
	return nil
}

func (w lintCyclomatic) Visit(n ast.Node) ast.Visitor { // MATCH /parameter 'n' seems to be unused, consider removing or renaming it as _/
	f := w.file
	for _, decl := range f.AST.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			c := complexity(fn)
			if c > w.complexity {
				w.onFailure(lint.Failure{
					Confidence: 1,
					Category:   "maintenance",
					Failure:    fmt.Sprintf("function %s has cyclomatic complexity %d", funcName(fn), c),
					Node:       fn,
				})
			}
		}
	}
	return nil
}

func ext۰time۰Sleep(fr *frame, args []value) value { // MATCH /parameter 'fr' seems to be unused, consider removing or renaming it as _/
	time.Sleep(time.Duration(args[0].(int64)))
	return nil
}
