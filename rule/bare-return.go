package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// BareReturnRule lints given else constructs.
type BareReturnRule struct{}

// Apply applies the rule to given file.
func (r *BareReturnRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintBareReturnRule{onFailure: onFailure}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (r *BareReturnRule) Name() string {
	return "bare-return"
}

type lintBareReturnRule struct {
	onFailure func(lint.Failure)
}

func (r lintBareReturnRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		r.checkFunc(n.Type.Results, n.Body)
	case *ast.FuncLit: // to cope with deferred functions and go-routines
		r.checkFunc(n.Type.Results, n.Body)
	}

	return r
}

// checkFunc will verify if the given function has named result and bare returns
func (r lintBareReturnRule) checkFunc(results *ast.FieldList, body *ast.BlockStmt) {
	hasNamedResults := results != nil && len(results.List) > 0 && results.List[0].Names != nil
	if !hasNamedResults || body == nil {
		return // nothing to do
	}

	brf := bareReturnFinder{r.onFailure}
	ast.Walk(brf, body)
}

type bareReturnFinder struct {
	onFailure func(lint.Failure)
}

func (r bareReturnFinder) Visit(node ast.Node) ast.Visitor {
	_, ok := node.(*ast.FuncLit)
	if ok {
		// skip analysing function literals
		// they will be analysed by the lintBareReturnRule.Visit method
		return nil
	}

	rs, ok := node.(*ast.ReturnStmt)
	if !ok {
		return r
	}

	if len(rs.Results) > 0 {
		return r
	}

	r.onFailure(lint.Failure{
		Confidence: 1,
		Node:       rs,
		Failure:    "avoid using bare returns, please add return expressions",
	})

	return r
}
