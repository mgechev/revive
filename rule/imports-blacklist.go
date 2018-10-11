package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ImportsBlacklistRule lints given else constructs.
type ImportsBlacklistRule struct{}

// Apply applies the rule to given file.
func (r *ImportsBlacklistRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	blacklist := make(map[string]bool, len(arguments))

	for _, arg := range arguments {
		argStr, ok := arg.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid argument to the imports-blacklist rule. Expecting a string, got %T", arg))
		}
		// we add quotes if nt present, because when parsed, the value of the AST node, will be quoted
		if len(argStr) > 2 && argStr[0] != '"' && argStr[len(argStr)-1] != '"' {
			argStr = fmt.Sprintf(`"%s"`, argStr)
		}
		blacklist[argStr] = true
	}

	fileAst := file.AST
	walker := blacklistedImports{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
		blacklist: blacklist,
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ImportsBlacklistRule) Name() string {
	return "imports-blacklist"
}

type blacklistedImports struct {
	file      *lint.File
	fileAst   *ast.File
	onFailure func(lint.Failure)
	blacklist map[string]bool
}

func (w blacklistedImports) Visit(_ ast.Node) ast.Visitor {
	for _, is := range w.fileAst.Imports {
		if is.Path != nil && !w.file.IsTest() && w.blacklist[is.Path.Value] {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("should not use the following blacklisted import: %s", is.Path.Value),
				Node:       is,
				Category:   "imports",
			})
		}
	}
	return nil
}
