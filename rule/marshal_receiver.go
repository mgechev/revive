package rule

import (
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// MarshalReceiverRule lints marshal/unmarshal methods with incorrect receiver types.
type MarshalReceiverRule struct{}

// Name returns the rule name.
func (*MarshalReceiverRule) Name() string {
	return "marshal-receiver"
}

// Apply applies the rule to given file.
func (*MarshalReceiverRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			continue
		}

		name := fn.Name.Name
		isMarshall := isMarshalMethod(name)
		isUnmarshall := !isMarshall && isUnmarshalMethod(name)
		if !(isMarshall || isUnmarshall) {
			continue
		}

		recv := fn.Recv.List[0]
		_, isPtr := recv.Type.(*ast.StarExpr)

		msg := ""
		switch {
		case isMarshal && isPtr:
			msg = " method should use a value receiver, not a pointer receiver"
		case isUnmarshal && !isPtr:
			msg = " method should use a pointer receiver, not a value receiver"
		default:
			continue // nothing to say about the method declaration
		}

		failures = append(failures, lint.Failure{
			Node:       decl,
			Confidence: 1,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    name + msg
		})
	}

	return failures
}

func isMarshalOrUnmarshalMethod(name string) bool {
	return isMarshalMethod(name) || isUnmarshalMethod(name)
}

func isMarshalMethod(name string) bool {
	return name == "MarshalJSON" || name == "MarshalYAML"
}

func isUnmarshalMethod(name string) bool {
	return name == "UnmarshalJSON" || name == "UnmarshalYAML"
}
