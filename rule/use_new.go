package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/internal/astutils"
	"github.com/mgechev/revive/lint"
)

// UseNewRule implements a rule that proposes using new(expr) when possible.
type UseNewRule struct{}

// Apply applies the rule to given file.
func (r *UseNewRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !file.Pkg.IsAtLeastGoVersion(lint.Go126) {
		return nil
	}

	var failures []lint.Failure
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}

		failures = append(failures, r.lintFunction(funcDecl)...)
	}

	return failures
}

// Name returns the rule name.
func (*UseNewRule) Name() string {
	return "use-new"
}

func (r *UseNewRule) lintFunction(funcDecl *ast.FuncDecl) []lint.Failure {
	if !r.isNewValueFunc(funcDecl) {
		return nil
	}

	return []lint.Failure{
		{
			Failure:    fmt.Sprintf(`calls to "%s(value)" can be replaced by a call to "new(value)"`, funcDecl.Name.Name),
			Confidence: 1,
			Node:       funcDecl,
			Category:   lint.FailureCategoryOptimization,
		},
	}
}

// isNewValueFunc checks if the function is of the form:
//
//	func foo(p Type) *Type {
//	  return &p
//	}
func (*UseNewRule) isNewValueFunc(funcDecl *ast.FuncDecl) bool {
	if funcDecl.Type.Results == nil || len(funcDecl.Type.Results.List) != 1 {
		return false // not one return value
	}

	if funcDecl.Type.Params == nil || len(funcDecl.Type.Params.List) != 1 {
		return false // not one parameter
	}

	if len(funcDecl.Body.List) != 1 {
		return false // not one statement
	}

	paramTypes := astutils.GetTypeNames(funcDecl.Type.Params)
	resultTypes := astutils.GetTypeNames(funcDecl.Type.Results)
	if "*"+paramTypes[0] != resultTypes[0] {
		return false // return type is not pointer to parameter type
	}

	retStmt, ok := funcDecl.Body.List[0].(*ast.ReturnStmt)
	if !ok || len(retStmt.Results) != 1 {
		return false // not a return statement with one result
	}

	returnExpr := astutils.GoFmt(retStmt.Results[0]) // TODO use more efficient way to retrieve the id

	return returnExpr == "&"+funcDecl.Type.Params.List[0].Names[0].Name
}
