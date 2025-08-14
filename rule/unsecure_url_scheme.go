package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/mgechev/revive/lint"
)

// UnsecureURLSchemeRule checks if a file contains string literals with unsecure URL schemes (for example: http://... in place of https://...).
type UnsecureURLSchemeRule struct{}

// Apply applied the rule to the given file.
func (r *UnsecureURLSchemeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if file.IsTest() {
		return nil // skip test files
	}

	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintUnsecureURLSchemeRule{
		onFailure: onFailure,
	}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*UnsecureURLSchemeRule) Name() string {
	return "unsecure-url-scheme"
}

type lintUnsecureURLSchemeRule struct {
	onFailure func(lint.Failure)
}

func (w lintUnsecureURLSchemeRule) Visit(node ast.Node) ast.Visitor {
	n, ok := node.(*ast.BasicLit)
	if !ok || n.Kind != token.STRING {
		return w // not a string litereal
	}

	value, _ := strconv.Unquote(n.Value)
	var scheme string
	switch {
	case strings.HasPrefix(value, `http://`):
		scheme = "http"
	case strings.HasPrefix(value, `ws://`):
		scheme = "ws"
	default:
		return nil // not an URL or not an unsecure one
	}

	if strings.Contains(value, "localhost") || strings.Contains(value, "127.0.0.1") {
		return nil // do not fail on local URL
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Failure:    fmt.Sprintf("preffer secure protocol %s over %s", scheme+"s", scheme),
		Node:       n,
	})

	return nil
}
