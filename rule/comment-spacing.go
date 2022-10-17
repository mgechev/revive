package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// CommentSpacings Rule check the whether there is a space between
// the comment symbol( // ) and the start of the comment text
type CommentSpacingsRule struct{}

func (*CommentSpacingsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	fileAST := file.AST
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	whiteList := map[string]bool{
		"//go:generate":              true,
		"//go:embed":                 true,
		"//go:build":                 true,
		"//go:linkname":              true,
		"//go:noescape":              true,
		"//go:noinline":              true,
		"//go:norace":                true,
		"//go:nowritebarrierrec":     true,
		"//go:nowritebarrier":        true,
		"//go:systemstack":           true,
		"//revive:disable":           true,
		"//revive:enable":            true,
		"//revive:ignore":            true,
		"//revive:ignore-next-line":  true,
		"//revive:ignore-next-lines": true,
		"//revive:ignore-line":       true,
		"//revive:ignore-lines":      true,
		"//revive:enable:exported":   true,
		"//revive:disable:exported":  true,
	}
	w := &lintCommentSpacings{whiteList: whiteList, onFailure: onFailure}
	ast.Walk(w, fileAST)
	return failures
}

func (*CommentSpacingsRule) Name() string {
	return "comment-spacings"
}

type lintCommentSpacings struct {
	whiteList map[string]bool
	onFailure func(lint.Failure)
}

func (l *lintCommentSpacings) Visit(n ast.Node) ast.Visitor {
	if node, ok := n.(*ast.CommentGroup); ok {
		for _, comment := range node.List {
			if _, ok := l.whiteList[comment.Text]; ok {
				continue
			}
			if comment.Text[0] == '/' && comment.Text[1] == '/' {
				if comment.Text[2] != ' ' {
					l.onFailure(lint.Failure{
						Node: n,
						Confidence: 1,
						Category:   "comments",
						Failure:    "no space between comment symbol and comment text",
					})
				}
			}
		}
	}
	return l
}
	
