package rule

import (
	"fmt"
	"go/ast"
	"go/types"
	"maps"

	"github.com/mgechev/revive/lint"
)

var defaultFilteredInterfaceNames = map[string]struct{}{}

var allFounded = map[string]struct{}{}

// ReturnsInterfaceTypeRule spots functions/methods returning an interface type.
type ReturnsInterfaceTypeRule struct {
	reportAll      bool     // enable/disable reporting of interface types
	searchingNames []string // set of user defined interface names used to filter results
}

// Apply applies the rule to given file.
func (r *ReturnsInterfaceTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if !r.reportAll {
		return failures
	}

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

			if !r.isFiltered(typ) {
				continue // configured to find this interface
			}

			returnName := r.returnFuncName(fn.Name.Name, signature)
			failures = append(failures, lint.Failure{
				Node:       fn,
				Confidence: 1.0,
				Failure:    returnName + " returns interface type " + typeStr,
			})
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
	r.reportAll = false
	r.searchingNames = []string{}
	allFounded = maps.Clone(defaultFilteredInterfaceNames)
	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument '%v' for '%s' rule. Expecting a k,v map, got %T", arguments[0], r.Name(), arguments[0])
	}

	for k, v := range args {
		if isRuleOption(k, "reportAll") {
			reportAll, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid argument '%v' fo '%s' rule; need bool but got %T`, k, r.Name(), k)
			}
			r.reportAll = reportAll
			continue
		}

		if !isRuleOption(k, "searchingNames") {
			return fmt.Errorf("invalid argument '%v' of '%s' rule configuration: searching-names expected. got '%v'", k, r.Name(), k)
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
			allFounded[name] = struct{}{}
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

// isFiltered helper function to check if type is filtered, decission is based
// on values from defaultFilteredNames(default values)+searchinNames(provided by user from config).
func (r *ReturnsInterfaceTypeRule) isFiltered(t types.Type) bool {
	if len(allFounded) == 0 {
		return true
	}
	name := r.getNameForType(t)
	_, filtered := allFounded[name]
	return name != "" && filtered
}

// getFilteredTypes helper function to get all filtered types.
func (*ReturnsInterfaceTypeRule) getFilteredTypes() map[string]struct{} {
	return allFounded
}

// DefaultFilteredTypes helper function to get by default filtered types.
func (*ReturnsInterfaceTypeRule) DefaultFilteredTypes() map[string]struct{} {
	return defaultFilteredInterfaceNames
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
