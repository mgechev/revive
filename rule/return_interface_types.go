package rule

import (
	"fmt"
	"go/ast"
	"go/types"
	"strconv"

	"github.com/mgechev/revive/lint"
)

// ReturnInterfaceTypesRule spot functions/methods returning an interface type.
type ReturnInterfaceTypesRule struct {
	stopOnFirst             bool         // stop on first founded interface type in function/method
	userDefinedIgnoredNames ignoredNames // set of user defined ignored interface names
}

// Apply applies the rule to given file.
func (r *ReturnInterfaceTypesRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
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

			if r.userDefinedIgnoredNames.isIgnored(typ.String()) {
				continue
			}

			if _, ok := typ.Underlying().(*types.Interface); ok {
				returnName := r.returnName(fn.Name.Name, signature)
				failures = append(failures, lint.Failure{
					Node:       fn,
					Confidence: 1.0,
					Failure:    returnName + " returns interface type: " + typ.String(),
				})
				if r.stopOnFirst {
					break
				}
			}
		}
	}
	return failures
}

// Name returns the rule name.
func (*ReturnInterfaceTypesRule) Name() string {
	return "return-interface-types"
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *ReturnInterfaceTypesRule) Configure(arguments lint.Arguments) error {
	r.stopOnFirst = false
	r.userDefinedIgnoredNames = ignoredNames{}
	if len(arguments) == 0 {
		return nil
	}

	args, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf("invalid argument to the return-interface-types rule. Expecting a k,v map, got %T", arguments[0])
	}

	for k, v := range args {
		if isRuleOption(k, "stopOnFirst") {
			stop, ok := v.(bool)
			if !ok {
				return fmt.Errorf("invalid argument to the return-interface-types rule, expecting bool value. Got '%v' (%T)", v, v)
			}
			r.stopOnFirst = stop
		}
		if !isRuleOption(k, "userDefinedIgnoredNames") {
			continue
		}
		names, ok := v.([]any)
		if !ok {
			return fmt.Errorf("invalid argument to the return-interface-types rule, []string expected. Got '%v' (%T)", v, v)
		}
		for _, value := range names {
			typ, ok := value.(string)
			if !ok {
				return fmt.Errorf("invalid argument to the return-interface-types rule, string expected. Got '%v' (%T)", value, value)
			}
			r.userDefinedIgnoredNames.add(typ)
		}
	}

	return nil
}

func (*ReturnInterfaceTypesRule) returnName(functionName string, signature *types.Signature) string {
	returnName := functionName

	if signature.Recv() != nil {
		recvType := signature.Recv().Type()
		returnName = recvType.String() + "." + returnName
	}

	return returnName
}

type ignoredNames map[string]struct{}

func (in ignoredNames) add(name string) {
	in[strconv.Quote(name)] = struct{}{}
}

func (in ignoredNames) isIgnored(name string) bool {
	_, ignored := in[strconv.Quote(name)]
	return ignored
}
