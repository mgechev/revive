package rule

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ImportsBlacklistRule lints given else constructs.
type ImportsBlacklistRule struct{}

// Apply applies the rule to given file.
func (r *ImportsBlacklistRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	blacklist := make([]string, len(arguments))

	for i, arg := range arguments {
		argStr, ok := arg.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid argument to the imports-blacklist rule. Expecting a string, got %T", arg))
		}
		argStr = strings.ToLower(strings.TrimSpace(argStr))
		if len(argStr) > 2 && argStr[0] != '"' && argStr[len(argStr)-1] != '"' {
			argStr = fmt.Sprintf(`"%s"`, argStr)
		}
		blacklist[i] = strings.ToLower(strings.TrimSpace(argStr))
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
	blacklist []string
}

func (w blacklistedImports) Visit(n ast.Node) ast.Visitor {
	for _, is := range w.fileAst.Imports {
		if is.Path != nil {
			for _, blacklisted := range w.blacklist {
				if strings.ToLower(is.Path.Value) == blacklisted && !w.file.IsTest() {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Failure:    fmt.Sprintf("should not use the following blacklisted import: %s", is.Path.Value),
						Node:       is,
						Category:   "imports",
					})
				}
			}
		}
	}
	return nil
}
