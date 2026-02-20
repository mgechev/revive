package rule

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
	"sync"

	gopackages "golang.org/x/tools/go/packages"

	"github.com/mgechev/revive/lint"
)

// defaultBadPackageNames is the list of "bad" package names from https://go.dev/wiki/CodeReviewComments#package-names
// and https://go.dev/blog/package-names#bad-package-names.
// The rule warns about the usage of any package name in this list if skipPackageNameChecks is false.
// Values in the list should be lowercased.
var defaultBadPackageNames = map[string]struct{}{
	"api":           {},
	"common":        {},
	"interface":     {},
	"interfaces":    {},
	"misc":          {},
	"miscellaneous": {},
	"shared":        {},
	"type":          {},
	"types":         {},
	"util":          {},
	"utilities":     {},
	"utils":         {},
}

// PackageNamingRule is a rule that checks package names.
type PackageNamingRule struct {
	skipPackageNameChecks             bool                // if true - disable check for meaningless and user-defined bad package names
	extraBadPackageNames              map[string]struct{} // inactive if skipPackageNameChecks is false
	skipPackageNameCollisionWithGoStd bool                // if true - disable checks for collisions with Go standard library package names

	pkgNameAlreadyChecked syncSet // set of packages names already checked
	// stdPackageNames holds the names of standard library packages excluding internal and vendor.
	// populated only if skipPackageNameCollisionWithGoStd is false.
	// E.g., `net/http` stored as `http`, `math/rand/v2` - `rand` etc.
	stdPackageNames map[string]struct{}
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *PackageNamingRule) Configure(arguments lint.Arguments) error {
	r.pkgNameAlreadyChecked = syncSet{elements: map[string]struct{}{}}

	if len(arguments) == 1 {
		args, ok := arguments[0].(map[string]any)
		if !ok {
			return fmt.Errorf("invalid argument to the package-naming rule: expecting a k,v map, but got %T", arguments[0])
		}

		for k, v := range args {
			switch {
			case isRuleOption(k, "skipPackageNameChecks"):
				r.skipPackageNameChecks = fmt.Sprint(v) == "true"
			case isRuleOption(k, "extraBadPackageNames"):
				extraBadPackageNames, ok := v.([]any)
				if !ok {
					return fmt.Errorf("invalid argument to the package-naming rule: expecting extraBadPackageNames of type slice of strings, but got %T", v)
				}
				for i, name := range extraBadPackageNames {
					if r.extraBadPackageNames == nil {
						r.extraBadPackageNames = map[string]struct{}{}
					}
					n, ok := name.(string)
					if !ok {
						return fmt.Errorf("invalid argument to the package-naming rule: expecting element %d of extraBadPackageNames to be a string, but got %v(%T)", i, name, name)
					}
					r.extraBadPackageNames[strings.ToLower(n)] = struct{}{}
				}
			case isRuleOption(k, "skipPackageNameCollisionWithGoStd"):
				r.skipPackageNameCollisionWithGoStd = fmt.Sprint(v) == "true"
			}
		}
	}

	if !r.skipPackageNameCollisionWithGoStd && r.stdPackageNames == nil {
		pkgs, err := gopackages.Load(nil, "std")
		if err != nil {
			return fmt.Errorf("load std packages: %w", err)
		}

		r.stdPackageNames = map[string]struct{}{}
		for _, pkg := range pkgs {
			if isInternalOrVendorPackage(pkg.PkgPath) {
				continue
			}
			r.stdPackageNames[pkg.Name] = struct{}{}
		}
	}

	return nil
}

// isInternalOrVendorPackage reports whether the path represents an internal or vendor directory.
//
// Borrowed and modified from
// https://github.com/golang/pkgsite/blob/84333735ffe124f7bd904805fd488b93841de49f/internal/postgres/search.go#L1009-L1016
func isInternalOrVendorPackage(path string) bool {
	for p := range strings.SplitSeq(path, "/") {
		if p == "internal" || p == "vendor" {
			return true
		}
	}
	return false
}

// Apply applies the rule to given file.
func (r *PackageNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	if r.skipPackageNameChecks {
		return failures
	}

	fileDir := filepath.Dir(file.Name)

	r.pkgNameAlreadyChecked.Lock()
	defer r.pkgNameAlreadyChecked.Unlock()
	if r.pkgNameAlreadyChecked.has(fileDir) {
		return failures
	}
	r.pkgNameAlreadyChecked.add(fileDir) // mark this package as already checked

	pkgNameNode := file.AST.Name
	pkgName := pkgNameNode.Name
	pkgNameLower := strings.ToLower(pkgName)

	// Check if top level package
	if pkgNameLower == "pkg" && filepath.Base(fileDir) != pkgName {
		onFailure(r.pkgNameFailure(pkgNameNode, "should not have a root level package called pkg"))
		return failures
	}

	if _, ok := r.extraBadPackageNames[pkgNameLower]; ok {
		onFailure(r.pkgNameFailure(pkgNameNode, "avoid bad package names"))
		return failures
	}

	if _, ok := defaultBadPackageNames[pkgNameLower]; ok {
		onFailure(r.pkgNameFailure(pkgNameNode, "avoid meaningless package names"))
		return failures
	}

	if !r.skipPackageNameCollisionWithGoStd {
		if _, ok := r.stdPackageNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "avoid package names that conflict with Go standard library package names"))
		}
	}

	// Package names need slightly different handling than other names.
	if strings.Contains(pkgName, "_") && !strings.HasSuffix(pkgName, "_test") {
		onFailure(r.pkgNameFailure(pkgNameNode, "don't use an underscore in package name"))
	}
	if hasUpperCaseLetter(pkgName) {
		onFailure(r.pkgNameFailure(pkgNameNode, "don't use MixedCaps in package names; %s should be %s", pkgName, pkgNameLower))
	}

	return failures
}

// Name returns the rule name.
func (*PackageNamingRule) Name() string {
	return "package-naming"
}

func (*PackageNamingRule) pkgNameFailure(node ast.Node, msg string, args ...any) lint.Failure {
	return lint.Failure{
		Failure:    fmt.Sprintf(msg, args...),
		Confidence: 1,
		Node:       node,
		Category:   lint.FailureCategoryNaming,
	}
}

type syncSet struct {
	sync.Mutex

	elements map[string]struct{}
}

func (sm *syncSet) has(s string) bool {
	_, result := sm.elements[s]
	return result
}

func (sm *syncSet) add(s string) {
	sm.elements[s] = struct{}{}
}
