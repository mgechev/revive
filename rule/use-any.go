package rule

import (
	"github.com/mgechev/revive/lint"
	"go/ast"
)

// UseAnyRule lints given else constructs.
type UseAnyRule struct{}

// Apply applies the rule to given file.
func (*UseAnyRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := lintUseAny{
		commentPositions: getCommentsPositions(file.AST.Comments),
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}
	fileAst := file.AST
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*UseAnyRule) Name() string {
	return "use-any"
}

type lintUseAny struct {
	commentPositions []int
	onFailure        func(lint.Failure)
}

func (w lintUseAny) Visit(n ast.Node) ast.Visitor {
	it, ok := n.(*ast.InterfaceType)
	if !ok {
		return w
	}

	if len(it.Methods.List) != 0 {
		return w // it is not an empty interface
	}

	for _, pos := range w.commentPositions {
		if pos > int(it.Pos()) && pos < int(it.End()) {
			return w // it is a comment inside the interface
		}
	}

	w.onFailure(lint.Failure{
		Node:       n,
		Confidence: 1,
		Category:   "naming",
		Failure:    "since GO 1.18 'interface{}' can be replaced by 'any'",
	})

	return w
}

func getCommentsPositions(commentGroups []*ast.CommentGroup) []int {
	result := []int{}
	for _, commentGroup := range commentGroups {
		for _, comment := range commentGroup.List {
			result = append(result, int(comment.Pos()))
		}
	}

	return result
}
