package rule

import (
	"fmt"
	"go/ast"

	"github.com/mgechev/revive/lint"
)

// ContextAsArgumentRule lints given else constructs.
type ContextAsArgumentRule struct {
}

// Apply applies the rule to given file.
func (r *ContextAsArgumentRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	var allowTypesBefore []string
	if len(args) != 1 {
		panic(fmt.Sprintf("Invalid argument to the context-as-argument rule. Expecting a single k,v map only, got %v args", len(args)))
	}
	argKV, ok := args[0].(map[string]interface{})
	if !ok {
		panic(fmt.Sprintf("Invalid argument to the context-as-argument rule. Expecting a k,v map, got %T", args[0]))
	}
	for k, v := range argKV {
		switch k {
		case "allowTypesBefore":
			typesBefore, ok := v.([]string)
			if !ok {
				panic(fmt.Sprintf("Invalid argument to the context-as-argument.allowTypesBefore rule. Expecting a []string, got %T", v))
			}
			allowTypesBefore = typesBefore
		default:
			panic(fmt.Sprintf("Invalid argument to the context-as-argument rule. Unrecognized key %s", k))
		}
	}

	var failures []lint.Failure

	fileAst := file.AST
	walker := lintContextArguments{
		allowTypesBefore: allowTypesBefore,
		file:             file,
		fileAst:          fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ContextAsArgumentRule) Name() string {
	return "context-as-argument"
}

type lintContextArguments struct {
	file             *lint.File
	fileAst          *ast.File
	allowTypesBefore []string
	onFailure        func(lint.Failure)
}

func (w lintContextArguments) Visit(n ast.Node) ast.Visitor {
	fn, ok := n.(*ast.FuncDecl)
	if !ok || len(fn.Type.Params.List) <= 1 {
		return w
	}
	allowTypesLookup := make(map[string]struct{}, len(w.allowTypesBefore))
	for _, v := range w.allowTypesBefore {
		allowTypesLookup[v] = struct{}{}
	}

	fnArgs := fn.Type.Params.List
	// trim off all types we've been configured to skip
	for len(fnArgs) > 0 {
		typeStr, ok := astExprTypeStr(fnArgs[0].Type)
		if !ok {
			// assume we're done. This can happen, for example, with a function type
			// argument (`func(x func()...)` which we choose not to represent.
			break
		}
		_, isAllowed := allowTypesLookup[typeStr]
		if isAllowed {
			// trim
			fnArgs = fnArgs[1:]
		} else {
			// first non-trimmed argument means we've trimmed all prefix args
			break
		}
	}

	if len(fnArgs) <= 1 {
		return w
	}
	// A context.Context should be the first parameter of a function.
	// Flag any that show up after the first.
	previousArgIsCtx := isPkgDot(fnArgs[0].Type, "context", "Context")
	for _, arg := range fnArgs[1:] {
		argIsCtx := isPkgDot(arg.Type, "context", "Context")
		if argIsCtx && !previousArgIsCtx {
			w.onFailure(lint.Failure{
				Node:       arg,
				Category:   "arg-order",
				Failure:    "context.Context should be the first parameter of a function",
				Confidence: 0.9,
			})
			break // only flag one
		}
		previousArgIsCtx = argIsCtx
	}
	return w
}
