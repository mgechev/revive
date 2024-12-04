package rule

import (
	"fmt"
	"go/ast"
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// UnusedReceiverRule lints unused receivers in functions.
type UnusedReceiverRule struct {
	// regex to check if some name is valid for unused parameter, "^_$" by default
	allowRegex *regexp.Regexp
	failureMsg string

	configureOnce sync.Once
}

func (r *UnusedReceiverRule) configure(args lint.Arguments) {
	// while by default args is an array, i think it's good to provide structures inside it by default, not arrays or primitives
	// it's more compatible to JSON nature of configurations
	r.allowRegex = allowBlankIdentifierRegex
	r.failureMsg = "method receiver '%s' is not referenced in method's body, consider removing or renaming it as _"
	if len(args) == 0 {
		return
	}
	// Arguments = [{}]
	options := args[0].(map[string]any)

	allowRegexParam, ok := options["allowRegex"]
	if !ok {
		return
	}
	// Arguments = [{allowRegex="^_"}]
	allowRegexStr, ok := allowRegexParam.(string)
	if !ok {
		panic(fmt.Errorf("error configuring [unused-receiver] rule: allowRegex is not string but [%T]", allowRegexParam))
	}
	var err error
	r.allowRegex, err = regexp.Compile(allowRegexStr)
	if err != nil {
		panic(fmt.Errorf("error configuring [unused-receiver] rule: allowRegex is not valid regex [%s]: %v", allowRegexStr, err))
	}
	r.failureMsg = "method receiver '%s' is not referenced in method's body, consider removing or renaming it to match " + r.allowRegex.String()
}

// Apply applies the rule to given file.
func (r *UnusedReceiverRule) Apply(file *lint.File, args lint.Arguments) []lint.Failure {
	r.configureOnce.Do(func() { r.configure(args) })
	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		isMethod := ok && funcDecl.Recv != nil
		if !isMethod {
			continue
		}

		rec := funcDecl.Recv.List[0] // safe to access only the first (unique) element of the list
		if len(rec.Names) < 1 {
			continue // the receiver is anonymous: func (aType) Foo(...) ...
		}

		recID := rec.Names[0]
		if recID.Name == "_" {
			continue // the receiver is already named _
		}

		if r.allowRegex != nil && r.allowRegex.FindStringIndex(recID.Name) != nil {
			continue
		}

		// inspect the func body looking for references to the receiver id
		selectReceiverUses := func(n ast.Node) bool {
			ident, isAnID := n.(*ast.Ident)

			return isAnID && ident.Obj == recID.Obj
		}
		receiverUses := pick(funcDecl.Body, selectReceiverUses)

		if len(receiverUses) > 0 {
			continue // the receiver is referenced in the func body
		}

		failures = append(failures, lint.Failure{
			Confidence: 1,
			Node:       recID,
			Category:   "bad practice",
			Failure:    fmt.Sprintf(r.failureMsg, recID.Name),
		})
	}

	return failures
}

// Name returns the rule name.
func (*UnusedReceiverRule) Name() string {
	return "unused-receiver"
}
