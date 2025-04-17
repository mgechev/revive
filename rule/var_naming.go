package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

var anyCapsRE = regexp.MustCompile(`[A-Z]`)

var allCapsRE = regexp.MustCompile(`^[A-Z0-9_]+$`)

// regexp for constant names like `SOME_CONST`, `SOME_CONST_2`, `X123_3`, `_SOME_PRIVATE_CONST` (#851, #865)
var upperCaseConstRE = regexp.MustCompile(`^_?[A-Z][A-Z\d]*(_[A-Z\d]+)*$`)

var knownNameExceptions = map[string]bool{
	"LastInsertId": true, // must match database/sql
	"kWh":          true,
}

// defaultBadPackageNames is the list of "bad" package names from https://go.dev/wiki/CodeReviewComments#package-names
// and https://go.dev/blog/package-names#bad-package-names.
// The rule warns about the usage of any package name in this list if skipPackageNameChecks is false.
// Values in the list should be lowercased.
var defaultBadPackageNames = map[string]struct{}{
	"common":     {},
	"interfaces": {},
	"misc":       {},
	"types":      {},
	"util":       {},
	"utils":      {},
}

// VarNamingRule lints the name of a variable.
type VarNamingRule struct {
	allowList             []string
	blockList             []string
	allowUpperCaseConst   bool                // if true - allows to use UPPER_SOME_NAMES for constants
	skipPackageNameChecks bool                // check for meaningless and user-defined bad package names
	extraBadPackageNames  map[string]struct{} // inactive if skipPackageNameChecks is false
	pkgNameAlreadyChecked set                 // set of packages names already checked
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *VarNamingRule) Configure(arguments lint.Arguments) error {
	r.pkgNameAlreadyChecked = set{elements: map[string]struct{}{}}

	if len(arguments) >= 1 {
		list, err := getList(arguments[0], "allowlist")
		if err != nil {
			return err
		}
		r.allowList = list
	}

	if len(arguments) >= 2 {
		list, err := getList(arguments[1], "blocklist")
		if err != nil {
			return err
		}
		r.blockList = list
	}

	if len(arguments) >= 3 {
		// not pretty code because should keep compatibility with TOML (no mixed array types) and new map parameters
		thirdArgument := arguments[2]
		asSlice, ok := thirdArgument.([]any)
		if !ok {
			return fmt.Errorf("invalid third argument to the var-naming rule. Expecting a %s of type slice, got %T", "options", arguments[2])
		}
		if len(asSlice) != 1 {
			return fmt.Errorf("invalid third argument to the var-naming rule. Expecting a %s of type slice, of len==1, but %d", "options", len(asSlice))
		}
		args, ok := asSlice[0].(map[string]any)
		if !ok {
			return fmt.Errorf("invalid third argument to the var-naming rule. Expecting a %s of type slice, of len==1, with map, but %T", "options", asSlice[0])
		}
		for k, v := range args {
			switch {
			case isRuleOption(k, "upperCaseConst"):
				r.allowUpperCaseConst = fmt.Sprint(v) == "true"
			case isRuleOption(k, "skipPackageNameChecks"):
				r.skipPackageNameChecks = fmt.Sprint(v) == "true"
			case isRuleOption(k, "extraBadPackageNames"):
				extraBadPackageNames, ok := v.([]string)
				if !ok {
					return fmt.Errorf("invalid third argument to the var-naming rule. Expecting extraBadPackageNames of type slice of strings, but %T", v)
				}
				for _, name := range extraBadPackageNames {
					if r.extraBadPackageNames == nil {
						r.extraBadPackageNames = map[string]struct{}{}
					}
					r.extraBadPackageNames[strings.ToLower(name)] = struct{}{}
				}
			}
		}
	}
	return nil
}

// Apply applies the rule to given file.
func (r *VarNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	r.applyPackageCheckRules(file, onFailure)

	fileAst := file.AST
	walker := lintNames{
		file:           file,
		fileAst:        fileAst,
		allowList:      r.allowList,
		blockList:      r.blockList,
		onFailure:      onFailure,
		upperCaseConst: r.allowUpperCaseConst,
	}

	ast.Walk(&walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*VarNamingRule) Name() string {
	return "var-naming"
}

func (r *VarNamingRule) applyPackageCheckRules(file *lint.File, onFailure func(failure lint.Failure)) {
	fileDir := path.Dir(file.Name)

	// Protect pkgsWithNameFailure from concurrent modifications
	r.pkgNameAlreadyChecked.Lock()
	defer r.pkgNameAlreadyChecked.Unlock()
	defer r.pkgNameAlreadyChecked.add(fileDir) // mark this package as already checked

	// Do we need to proceed with the checks on package name?
	mustSkipCheck := r.pkgNameAlreadyChecked.has(fileDir) || r.skipPackageNameChecks
	if mustSkipCheck {
		return
	}

	pkgNameNode := file.AST.Name
	pkgName := pkgNameNode.Name
	pkgNameLower := strings.ToLower(pkgName)
	if _, ok := r.extraBadPackageNames[pkgNameLower]; ok {
		onFailure(r.pkgNameFailure(pkgNameNode, "avoid bad package names"))
		return
	}

	if _, ok := defaultBadPackageNames[pkgNameLower]; ok {
		onFailure(r.pkgNameFailure(pkgNameNode, "avoid meaningless package names"))
		return
	}

	// Package names need slightly different handling than other names.
	if strings.Contains(pkgName, "_") && !strings.HasSuffix(pkgName, "_test") {
		onFailure(r.pkgNameFailure(pkgNameNode, "don't use underscores in package names"))
	}
	if anyCapsRE.MatchString(pkgName) {
		onFailure(r.pkgNameFailure(pkgNameNode, "don't use MixedCaps in package names; %s should be %s", pkgName, pkgNameLower))
	}
}

func (*VarNamingRule) pkgNameFailure(node ast.Node, msg string, args ...any) lint.Failure {
	return lint.Failure{
		Failure:    fmt.Sprintf(msg, args...),
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryNaming,
	}
}

func (w *lintNames) checkList(fl *ast.FieldList, thing string) {
	if fl == nil {
		return
	}
	for _, f := range fl.List {
		for _, id := range f.Names {
			w.check(id, thing)
		}
	}
}

func (w *lintNames) check(id *ast.Ident, thing string) {
	if id.Name == "_" {
		return
	}
	if knownNameExceptions[id.Name] {
		return
	}

	// #851 upperCaseConst support
	// if it's const
	if thing == token.CONST.String() && w.upperCaseConst && upperCaseConstRE.MatchString(id.Name) {
		return
	}

	// Handle two common styles from other languages that don't belong in Go.
	if len(id.Name) >= 5 && allCapsRE.MatchString(id.Name) && strings.Contains(id.Name, "_") {
		w.onFailure(lint.Failure{
			Failure:    "don't use ALL_CAPS in Go names; use CamelCase",
			Confidence: 0.8,
			Node:       id,
			Category:   lint.FailureCategoryNaming,
		})
		return
	}

	should := lint.Name(id.Name, w.allowList, w.blockList)
	if id.Name == should {
		return
	}

	if len(id.Name) > 2 && strings.Contains(id.Name[1:], "_") {
		w.onFailure(lint.Failure{
			Failure:    fmt.Sprintf("don't use underscores in Go names; %s %s should be %s", thing, id.Name, should),
			Confidence: 0.9,
			Node:       id,
			Category:   lint.FailureCategoryNaming,
		})
		return
	}
	w.onFailure(lint.Failure{
		Failure:    fmt.Sprintf("%s %s should be %s", thing, id.Name, should),
		Confidence: 0.8,
		Node:       id,
		Category:   lint.FailureCategoryNaming,
	})
}

type lintNames struct {
	file           *lint.File
	fileAst        *ast.File
	onFailure      func(lint.Failure)
	allowList      []string
	blockList      []string
	upperCaseConst bool
}

func (w *lintNames) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.AssignStmt:
		if v.Tok == token.ASSIGN {
			return w
		}
		for _, exp := range v.Lhs {
			if id, ok := exp.(*ast.Ident); ok {
				w.check(id, "var")
			}
		}
	case *ast.FuncDecl:
		funcName := v.Name.Name
		if w.file.IsTest() &&
			(strings.HasPrefix(funcName, "Example") ||
				strings.HasPrefix(funcName, "Test") ||
				strings.HasPrefix(funcName, "Benchmark") ||
				strings.HasPrefix(funcName, "Fuzz")) {
			return w
		}

		thing := "func"
		if v.Recv != nil {
			thing = "method"
		}

		// Exclude naming warnings for functions that are exported to C but
		// not exported in the Go API.
		// See https://github.com/golang/lint/issues/144.
		if ast.IsExported(v.Name.Name) || !isCgoExported(v) {
			w.check(v.Name, thing)
		}

		w.checkList(v.Type.Params, thing+" parameter")
		w.checkList(v.Type.Results, thing+" result")
	case *ast.GenDecl:
		if v.Tok == token.IMPORT {
			return w
		}

		thing := v.Tok.String()
		for _, spec := range v.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				w.check(s.Name, thing)
			case *ast.ValueSpec:
				for _, id := range s.Names {
					w.check(id, thing)
				}
			}
		}
	case *ast.InterfaceType:
		// Do not check interface method names.
		// They are often constrained by the method names of concrete types.
		for _, x := range v.Methods.List {
			ft, ok := x.Type.(*ast.FuncType)
			if !ok { // might be an embedded interface name
				continue
			}
			w.checkList(ft.Params, "interface method parameter")
			w.checkList(ft.Results, "interface method result")
		}
	case *ast.RangeStmt:
		if v.Tok == token.ASSIGN {
			return w
		}
		if id, ok := v.Key.(*ast.Ident); ok {
			w.check(id, "range var")
		}
		if id, ok := v.Value.(*ast.Ident); ok {
			w.check(id, "range var")
		}
	case *ast.StructType:
		for _, f := range v.Fields.List {
			for _, id := range f.Names {
				w.check(id, "struct field")
			}
		}
	}
	return w
}

func getList(arg any, argName string) ([]string, error) {
	args, ok := arg.([]any)
	if !ok {
		return nil, fmt.Errorf("invalid argument to the var-naming rule. Expecting a %s of type slice with initialisms, got %T", argName, arg)
	}
	var list []string
	for _, v := range args {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("invalid %v values of the var-naming rule. Expecting slice of strings but got element of type %T", v, arg)
		}
		list = append(list, val)
	}
	return list, nil
}

type set struct {
	sync.Mutex
	elements map[string]struct{}
}

func (sm *set) has(s string) bool {
	_, result := sm.elements[s]
	return result
}

func (sm *set) add(s string) {
	sm.elements[s] = struct{}{}
}
