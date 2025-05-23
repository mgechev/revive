package rule

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/mgechev/revive/lint"
)

// UnnecessaryFormatRule spots calls to formatting functions without leveraging formatting directives.
type UnnecessaryFormatRule struct{}

// Apply applies the rule to given file.
func (*UnnecessaryFormatRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	walker := lintUnnecessaryFormat{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*UnnecessaryFormatRule) Name() string {
	return "unnecessary-format"
}

type lintUnnecessaryFormat struct {
	file      *lint.File
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

type formattingSpec struct {
	formatArgPosition byte
	alternative       string
}

var formattingFuncs = map[string]formattingSpec{
	"fmt.Appendf":   {1, "fmt.Append"},
	"fmt.Errorf":    {0, "errors.New"},
	"fmt.Fprintf":   {1, "fmt.Fprint"},
	"fmt.Fscanf":    {1, "fmt.Fscan or fmt.Fscanln"},
	"fmt.Printf":    {0, "fmt.Print or fmt.Println"},
	"fmt.Scanf":     {0, "fmt.Scan"},
	"fmt.Sprintf":   {0, "fmt.Sprint or just the string itself"},
	"fmt.Sscanf":    {1, "fmt.Sscan"},
	"log.Fatalf":    {0, "log.Fatal"},
	"log.Panicf":    {0, "log.Panic"},
	"log.Printf":    {0, "log.Print"},
	"logger.Fatalf": {0, "logger.Fatal"},
	"logger.Panicf": {0, "logger.Panic"},
	"logger.Printf": {0, "logger.Print"},
}

func (w lintUnnecessaryFormat) Visit(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if !ok || len(ce.Args) < 1 {
		return w
	}

	funcName := gofmt(ce.Fun)
	spec, ok := formattingFuncs[funcName]
	if !ok {
		return w
	}

	pos := int(spec.formatArgPosition)
	if len(ce.Args) <= pos {
		return w // not enough params /!\
	}

	format := gofmt(ce.Args[pos])

	if !(format[0] == '"') || strings.Contains(format, `%`) {
		return w
	}

	failure := lint.Failure{
		Category:   lint.FailureCategoryOptimization,
		Node:       ce.Fun,
		Confidence: 1,
		Failure:    fmt.Sprintf("unnecessary use of formatting function %s, you can replace it with %s", funcName, spec.alternative),
	}

	w.onFailure(failure)

	return w
}
