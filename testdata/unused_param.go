package fixtures

import (
	"fmt"
	"go/ast"
	"os"
	"runtime"
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
	a := param.field
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

func f6(unused string) { // MATCH /parameter 'unused' seems to be unused, consider removing or renaming it as _/
	switch unused := runtime.GOOS; unused {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.", unused)
	}
	for unused := 0; unused < 10; unused++ {
		sum += unused
	}
	{
		unused := 1
	}
}

func f6bis(unused string) {
	switch unused := runtime.GOOS; unused {
	case "darwin":
		fmt.Println("OS X.")
	case "linux":
		fmt.Println("Linux.")
	default:
		fmt.Printf("%s.", unused)
	}
	for unused := 0; unused < 10; unused++ {
		sum += unused
	}
	{
		unused := 1
	}

	fmt.Print(unused)
}

func f7(pl int) {
	for i := 0; pl < i; i-- {

	}
}

func getCompareFailCause(n *node, which int, prevValue string, prevIndex uint64) string {
	switch which {
	case CompareIndexNotMatch:
		return fmt.Sprintf("[%v != %v]", prevIndex, n.ModifiedIndex)
	case CompareValueNotMatch:
		return fmt.Sprintf("[%v != %v]", prevValue, n.Value)
	default:
		return fmt.Sprintf("[%v != %v] [%v != %v]", prevValue, n.Value, prevIndex, n.ModifiedIndex)
	}
}

func assertSuccess(t *testing.T, baseDir string, fi os.FileInfo, src []byte, rules []lint.Rule, config map[string]lint.RuleConfig) error { // MATCH /parameter 'src' seems to be unused, consider removing or renaming it as _/
	l := lint.New(func(file string) ([]byte, error) {
		return os.ReadFile(baseDir + file)
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

func (c *chanList) remove(id uint32) {
	id -= c.offset
}

func (c *chanList) remove1(id uint32) {
	id *= c.offset
}

func (c *chanList) remove2(id uint32) {
	id /= c.offset
}

func (c *chanList) remove3(id uint32) {
	id += c.offset
}

func encodeFixed64Rpc(dAtA []byte, offset int, v uint64, i int) int {
	dAtA[offset+i] = uint8(v)

	return 8
}

func innerAnonymousFunctionWithoutUsage() {
	innerFunc := func(a int) {} // MATCH /parameter 'a' seems to be unused, consider removing or renaming it as _/
	innerFunc(1)
}

func innerAnonymousFunctionWithUsage() {
	innerFunc := func(a int) {
		a += 1
	}
	innerFunc(1)

	return someFunc(func(values []int) float64 { // MATCH /parameter 'values' seems to be unused, consider removing or renaming it as _/
		return 1.1
	})
}
