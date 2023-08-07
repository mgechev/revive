package rule

import (
	"fmt"
	"go/ast"
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// UnusedParamRule lints unused params in functions.
type UnusedParamRule struct {
	configured bool
	// regex to check if some name is valid for unused parameter, "^_$" by default
	allowRegex *regexp.Regexp
	sync.Mutex
}

func (r *UnusedParamRule) configure(args lint.Arguments) {
	// optimistic pre-check
	if r.configured {
		return
	}
	r.Lock()
	defer r.Unlock()
	if r.configured {
		return
	}
	r.configured = true
	// while by default args is an array, i think it's good to provide structures inside it by default, not arrays or primitives
	// it's more compatible to JSON nature of configurations
	if len(args) == 0 {
		return
	}
	// Arguments = [{}]
	options := args[0].(map[string]interface{})
	// Arguments = [{allowedRegex="^_"}]

	if allowedRegexParam, ok := options["allowRegex"]; ok {
		allowedRegexStr, ok := allowedRegexParam.(string)
		if !ok {
			panic(fmt.Errorf("error configuring [unused-parameter] rule: allowedRegex is not string but [%T]", allowedRegexParam))
		}
		var err error
		r.allowRegex, err = regexp.Compile(allowedRegexStr)
		if err != nil {
			panic(fmt.Errorf("error configuring [unused-parameter] rule: allowedRegex is not valid regex [%s]: %v", allowedRegexStr, err))
		}
	}
}

// Apply applies the rule to given file.
func (r *UnusedParamRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	r.configure(args)
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := lintUnusedParamRule{onFailure: onFailure, allowRegex: r.allowRegex}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*UnusedParamRule) Name() string {
	return "unused-parameter"
}

type lintUnusedParamRule struct {
	onFailure  func(lint.Failure)
	allowRegex *regexp.Regexp
}

func (w lintUnusedParamRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		params := retrieveNamedParams(n.Type.Params)
		if len(params) < 1 {
			return nil // skip, func without parameters
		}

		if n.Body == nil {
			return nil // skip, is a function prototype
		}

		// inspect the func body looking for references to parameters
		fselect := func(n ast.Node) bool {
			ident, isAnID := n.(*ast.Ident)

			if !isAnID {
				return false
			}

			_, isAParam := params[ident.Obj]
			if isAParam {
				params[ident.Obj] = false // mark as used
			}

			return false
		}
		_ = pick(n.Body, fselect, nil)

		for _, p := range n.Type.Params.List {
			for _, n := range p.Names {
				// проверка на соответствие паттерну допустимых не используемых параметров
				if w.allowRegex != nil && w.allowRegex.FindStringIndex(n.Name) != nil {
					continue
				}
				if params[n.Obj] {
					w.onFailure(lint.Failure{
						Confidence: 1,
						Node:       n,
						Category:   "bad practice",
						Failure:    fmt.Sprintf("parameter '%s' seems to be unused, consider removing or renaming it as _", n.Name),
					})
				}
			}
		}

		return nil // full method body already inspected
	}

	return w
}

func retrieveNamedParams(params *ast.FieldList) map[*ast.Object]bool {
	result := map[*ast.Object]bool{}
	if params.List == nil {
		return result
	}

	for _, p := range params.List {
		for _, n := range p.Names {
			if n.Name == "_" {
				continue
			}

			result[n.Obj] = true
		}
	}

	return result
}
