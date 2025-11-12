package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ConfusingEpochRule lints epoch time declarations.
type ConfusingEpochRule struct{}

// Apply applies the rule to given file.
func (*ConfusingEpochRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	walker := confusingEpoch{
		file: file,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck()
	ast.Walk(walker, file.AST)

	return failures
}

// Name returns the rule name.
func (*ConfusingEpochRule) Name() string {
	return "confusing-epoch"
}

type confusingEpoch struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

var epochUnits = map[string][]string{
	"Unix":      {"Sec", "Second", "Seconds"},
	"UnixMilli": {"Milli", "Ms"},
	"UnixMicro": {"Micro", "Microsecond", "Microseconds", "Us"},
	"UnixNano":  {"Nano", "Ns"},
}

func (w confusingEpoch) Visit(node ast.Node) ast.Visitor {
	switch v := node.(type) {
	case *ast.ValueSpec:
		// Handle var declarations
		for i, name := range v.Names {
			// Skip if no initialization value
			if i >= len(v.Values) {
				continue
			}

			w.check(name, v.Values[i])
		}
	case *ast.AssignStmt:
		// Handle both short variable declarations (:=) and regular assignments (=)
		if v.Tok != token.DEFINE && v.Tok != token.ASSIGN {
			return w
		}

		for i, lhs := range v.Lhs {
			if i >= len(v.Rhs) {
				continue
			}
			if ident, ok := lhs.(*ast.Ident); !ok || ident.Name == "_" {
					continue
			}
			
			w.check(ident, v.Rhs[i])
		}
	}

	return w
}

func (w confusingEpoch) check(name *ast.Ident, value ast.Expr) {
	call, ok := value.(*ast.CallExpr)
	if !ok {
		return
	}

	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return
	}

	// Check if the receiver is of type time.Time
	receiverType := w.file.Pkg.TypeOf(selector.X)
	if !isTime(receiverType) {
		return
	}

	methodName := selector.Sel.Name
	suffixes, ok := epochUnits[methodName]
	if !ok {
		return
	}

	varName := name.Name
	if !hasAnySuffix(varName, suffixes) {
		displaySuffixes := epochUnits[methodName]
		w.onFailure(lint.Failure{
			Confidence: 0.9,
			Node:       name,
			Category:   lint.FailureCategoryTime,
			Failure:    fmt.Sprintf("variable '%s' initialized with %s() should have a name end with one of %v", varName, methodName, displaySuffixes),
		})
	}
}

func isTime(typ types.Type) bool {
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

    obj := named.Obj()
    pkg := obj.Pkg()
	return pkg != nil && pkg.Path() == "time" && obj.Name() == "Time"
}

func hasAnySuffix(s string, suffixes []string) bool {
	lowerName := strings.ToLower(s)
	for _, suffix := range suffixes {
		if strings.HasSuffix(lowerName, strings.ToLower(suffix)) {
			return true
		}
	}
	return false
}
