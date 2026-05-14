package rule

import (
	"fmt"
	"go/ast"
	"go/types"

	"github.com/mgechev/revive/lint"
)

var ignoredInterfaceNames = map[string]struct{}{
	"error":       {},
	"any":         {},
	"interface{}": {},
}

// ReturnsInterfaceTypeRule spots functions/methods returning an interface type.
type ReturnsInterfaceTypeRule struct {
	stopOnFirst  bool                // stop on first found interface type in function/method
	ignoredNames map[string]struct{} // set of user defined ignored interface names
}

// Apply applies the rule to given file.
func (r *ReturnsInterfaceTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if err := file.Pkg.TypeCheck(); err != nil {
		return []lint.Failure{
			lint.NewInternalFailure(fmt.Sprintf("Unable to type check file %q: %v", file.Name, err)),
		}
	}

	info := file.Pkg.TypesInfo()

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)

		if !ok || fn.Name == nil {
			continue
		}

		obj := info.Defs[fn.Name]
		if obj == nil {
			continue
		}

		funcObj, ok := obj.(*types.Func)
		if !ok {
			continue
		}
		signature := funcObj.Type().(*types.Signature)

		results := signature.Results()

		for res := range results.Variables() {
			typ := res.Type()

			if _, ok := typ.Underlying().(*types.Interface); !ok {
				continue // not an interface
			}

			typeStr := typ.String()

			if r.isIgnored(typ) {
				continue // configured to ignore this interface
			}

			returnName := r.returnFuncName(fn.Name.Name, signature)
			failures = append(failures, lint.Failure{
				Node:       fn,
				Confidence: 1.0,
				Failure:    returnName + " returns interface type " + typeStr,
			})
			if r.stopOnFirst {
				break
			}
		}
	}
	return failures
}

// Name returns the rule name.
func (*ReturnsInterfaceTypeRule) Name() string {
	return "returns-interface-type"
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *ReturnsInterfaceTypeRule) Configure(arguments lint.Arguments) error {
	r.stopOnFirst = false
	r.ignoredNames = map[string]struct{}{}
	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument '%v' for '%s' rule. Expecting a k,v map, got %T", arguments[0], r.Name(), arguments[0])
	}

	for k, v := range args {
		if isRuleOption(k, "stopOnFirst") {
			stop, ok := v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument '%v' for '%s' rule, expecting bool value. got '%v' (%T)", k, r.Name(), v, v)
			}
			r.stopOnFirst = stop
			continue
		}
		if !isRuleOption(k, "ignoredNames") {
			continue
		}
		names, ok := v.([]any)
		if !ok {
			return fmt.Errorf("invalid format for entry '%v' of '%s' rule configuration: []string expected. got '%v' (%T)", k, r.Name(), v, v)
		}
		for _, p := range names {
			name, ok := p.(string)
			if !ok {
				return fmt.Errorf("invalid format for value in '%v' of '%s' rule configuration: string expected. got '%v' (%T)", k, r.Name(), p, p)
			}
			ignoredInterfaceNames[name] = struct{}{}
		}
	}
	return nil
}

// returnFuncName helper function to return name (with package name if defined).
func (*ReturnsInterfaceTypeRule) returnFuncName(functionName string, signature *types.Signature) string {
	returnName := functionName

	if signature.Recv() != nil {
		recvType := signature.Recv().Type()
		returnName = recvType.String() + "." + returnName
	}

	return returnName
}

// isIgnored helper function to check if type is ignored, decission is based
// on values from ignoredInterfaceNames(default values)+ignoredNames(provided by user from config).
func (r *ReturnsInterfaceTypeRule) isIgnored(t types.Type) bool {
	name := r.getNameForType(t)
	_, ignored := ignoredInterfaceNames[name]
	return name != "" && ignored
}

// getIgnoredTypes helper function to get all ignored types.
func (*ReturnsInterfaceTypeRule) getIgnoredTypes() map[string]struct{} {
	return ignoredInterfaceNames
}

// getNameForType helper function to get name from type
// if empty string that means we dont know what type it is,
// and it should be processed by rule and return proper warning.
func (*ReturnsInterfaceTypeRule) getNameForType(t types.Type) string {
	switch t := t.(type) {
	case *types.Named:
		return t.String()
	case *types.Alias:
		return t.String()
	case *types.Interface:
		return t.String()
	}
	return ""
}
