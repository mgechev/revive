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

// defaultBadNames is the list of "bad" package names from https://go.dev/blog/package-names#bad-package-names.
var defaultBadNames = map[string]struct{}{
	"common":     {},
	"interface":  {},
	"interfaces": {},
	"misc":       {},
	"type":       {},
	"types":      {},
	"util":       {},
	"utils":      {},
}

// extraBadNames is the list of additional "bad" package names that are not recommended.
var extraBadNames = map[string]struct{}{
	"api":           {},
	"helpers":       {},
	"miscellaneous": {},
	"models":        {},
	"shared":        {},
	"utilities":     {},
}

// commonStdNames is the list of standard library package names that are commonly used in Go programs.
var commonStdNames = map[string]string{
	"atomic":   "sync/atomic",
	"bytes":    "bytes",
	"context":  "context",
	"crypto":   "crypto",
	"errors":   "errors",
	"filepath": "path/filepath",
	"fmt":      "fmt",
	"http":     "net/http",
	"io":       "io",
	"json":     "encoding/json",
	"log":      "log",
	"maps":     "maps",
	"math":     "math",
	"net":      "net",
	"os":       "os",
	"path":     "path",
	"reflect":  "reflect",
	"regexp":   "regexp",
	"runtime":  "runtime",
	"slices":   "slices",
	"slog":     "log/slog",
	"sort":     "sort",
	"strings":  "strings",
	"sync":     "sync",
	"testing":  "testing",
	"time":     "time",
	"url":      "net/url",
	"xml":      "encoding/xml",
}

// PackageNamingRule is a rule that checks package names.
type PackageNamingRule struct {
	skipConventionChecks bool // if true - skip checks for package name conventions (e.g., no underscores, no MixedCaps etc.)
	skipTopLevelChecks   bool // if true - skip checks for top level package names (e.g., "pkg")

	skipDefaultBadNameChecks bool                // if true - enable check for default bad package names (e.g., "util", "misc" etc.)
	checkExtraBadNames       bool                // if true - enable check for extra bad package names (e.g., "helpers", "models" etc.)
	userDefinedBadNames      map[string]struct{} // set of user defined bad package names

	skipCollisionWithCommonStd bool // if true - skip checks for collisions with common Go standard library package names (e.g., "http", "json", "rand" etc.)

	checkCollisionWithAllStd bool // if true - enable checks for collisions with all Go standard library package names (including "version", "metrics" etc.)
	// allStdNames holds the names of standard library packages excluding internal and vendor.
	// Populated only if checkCollisionWithAllStd is true. `net/http` stored as `http`, `math/rand/v2` as `rand` etc.
	allStdNames map[string]string // name -> path

	alreadyCheckedNames syncSet // set of packages names already checked
}

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *PackageNamingRule) Configure(arguments lint.Arguments) error {
	r.alreadyCheckedNames = syncSet{elements: map[string]struct{}{}}

	if len(arguments) == 1 {
		args, ok := arguments[0].(map[string]any)
		if !ok {
			return fmt.Errorf("invalid argument to the package-naming rule: expecting a k,v map, but got %T", arguments[0])
		}

		for k, v := range args {
			switch {
			case isRuleOption(k, "skipConventionChecks"):
				r.skipConventionChecks = fmt.Sprint(v) == "true"
			case isRuleOption(k, "skipTopLevelChecks"):
				r.skipTopLevelChecks = fmt.Sprint(v) == "true"
			case isRuleOption(k, "skipDefaultBadNameChecks"):
				r.skipDefaultBadNameChecks = fmt.Sprint(v) == "true"
			case isRuleOption(k, "userDefinedBadNames"):
				userDefinedBadNames, ok := v.([]any)
				if !ok {
					return fmt.Errorf("invalid argument to the package-naming rule: expecting userDefinedBadNames of type slice of strings, but got %T", v)
				}
				for i, name := range userDefinedBadNames {
					if r.userDefinedBadNames == nil {
						r.userDefinedBadNames = map[string]struct{}{}
					}
					n, ok := name.(string)
					if !ok {
						return fmt.Errorf("invalid argument to the package-naming rule: expecting element %d of userDefinedBadNames to be a string, but got %v(%T)", i, name, name)
					}
					if n == "" {
						return fmt.Errorf("invalid argument to the package-naming rule: userDefinedBadNames cannot contain empty string (index %d)", i)
					}
					r.userDefinedBadNames[strings.ToLower(n)] = struct{}{}
				}
			case isRuleOption(k, "skipCollisionWithCommonStd"):
				r.skipCollisionWithCommonStd = fmt.Sprint(v) == "true"
			case isRuleOption(k, "checkCollisionWithAllStd"):
				r.checkCollisionWithAllStd = fmt.Sprint(v) == "true"
			}
		}
	}

	if r.checkCollisionWithAllStd && r.allStdNames == nil {
		pkgs, err := gopackages.Load(nil, "std")
		if err != nil {
			return fmt.Errorf("load std packages: %w", err)
		}

		r.allStdNames = map[string]string{}
		for _, pkg := range pkgs {
			if isInternalOrVendorPackage(pkg.PkgPath) {
				continue
			}
			r.allStdNames[pkg.Name] = pkg.PkgPath
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

	fileDir := filepath.Dir(file.Name)

	r.alreadyCheckedNames.Lock()
	defer r.alreadyCheckedNames.Unlock()
	if r.alreadyCheckedNames.has(fileDir) {
		return failures
	}
	r.alreadyCheckedNames.add(fileDir) // mark this package as already checked

	pkgNameNode := file.AST.Name
	pkgName := pkgNameNode.Name
	pkgNameLower := strings.ToLower(pkgName)

	if !r.skipConventionChecks {
		// Package names need slightly different handling than other names.
		if strings.Contains(pkgName, "_") && !strings.HasSuffix(pkgName, "_test") {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use package name %q that contains an underscore", pkgName))
		}
		if hasUpperCaseLetter(pkgName) {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use package name %q that contains MixedCaps", pkgName))
		}
	}

	if !r.skipTopLevelChecks {
		// Check if top level package
		if pkgNameLower == "pkg" && filepath.Base(fileDir) != pkgName {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q as a root level package name", pkgName))
			return failures
		}
	}

	if !r.skipDefaultBadNameChecks {
		if _, ok := defaultBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q because it is a bad package name according to https://go.dev/blog/package-names#bad-package-names", pkgName))
			return failures
		}
	}

	if r.checkExtraBadNames {
		if _, ok := extraBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q because it is a bad package name (extra)", pkgName))
			return failures
		}
	}

	if r.userDefinedBadNames != nil {
		if _, ok := r.userDefinedBadNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q because it is a bad package name (user-defined)", pkgName))
			return failures
		}
	}

	if !r.skipCollisionWithCommonStd {
		if std, ok := commonStdNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q because it conflicts with common Go standard library package %q", pkgName, std))
		}
	}

	if r.checkCollisionWithAllStd {
		if std, ok := r.allStdNames[pkgNameLower]; ok {
			onFailure(r.pkgNameFailure(pkgNameNode, "don't use %q because it conflicts with Go standard library package %q", pkgName, std))
		}
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
