package rule

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ConfusingNamingRule lints method names that differ only by capitalization
type ConfusingNamingRule struct{}

// Apply applies the rule to given file.
func (r *ConfusingNamingRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	walker := lintConfusingNames{
		methodNames: make(map[string][]*ast.Ident),
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(&walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ConfusingNamingRule) Name() string {
	return "confusing-naming"
}

//checkMethodName checks if a given method/function name is similar (just case differences) to other method/function of the same struct/file.
func checkMethodName(holder string, id *ast.Ident, w *lintConfusingNames) {
	name := strings.ToUpper(id.Name)
	if w.methodNames[holder] != nil {
		blackList := w.methodNames[holder]
		for _, n := range blackList {
			if strings.ToUpper(n.Name) == name {
				// confusing names
				w.onFailure(lint.Failure{
					Failure:    fmt.Sprintf("Method '%s' differs only by capitalization to method '%s'", id.Name, n.Name),
					Confidence: 1,
					Node:       id,
					Category:   "naming",
					URL:        "#TODO",
				})

				return
			}
		}
	}
	// update the black list
	w.methodNames[holder] = append(w.methodNames[holder], id)
}

type lintConfusingNames struct {
	methodNames map[string][]*ast.Ident // a map from struct names to method id nodes
	onFailure   func(lint.Failure)
}

const defaultStructName = "_" // used to map functions

//getStructName of a function receiver. Defaults to defaultStructName
func getStructName(r *ast.FieldList) string {
	result := defaultStructName

	if r == nil || len(r.List) < 1 {
		return result
	}

	t := r.List[0].Type

	if p, _ := t.(*ast.StarExpr); p != nil { // if a pointer receiver => dereference pointer receiver types
		t = p.X
	}

	if p, _ := t.(*ast.Ident); p != nil {
		result = p.Name
	}

	return result
}

func (w *lintConfusingNames) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.FuncDecl:
		// Exclude naming warnings for functions that are exported to C but
		// not exported in the Go API.
		// See https://github.com/golang/lint/issues/144.
		if ast.IsExported(v.Name.Name) || !isCgoExported(v) {
			checkMethodName(getStructName(v.Recv), v.Name, w)
		}
	default:
		// will add other checks like field names, struct names, etc.
	}

	return w
}
