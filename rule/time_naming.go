package rule

import (
	"fmt"
	"go/ast"
	"go/types"
	"log/slog"
	"strings"

	"github.com/mgechev/revive/lint"
)

// TimeNamingRule lints the name of a time variable.
type TimeNamingRule struct {
	logger *slog.Logger
}

// Apply applies the rule to given file.
func (r *TimeNamingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeNames{file, onFailure}

	if err := file.Pkg.TypeCheck(); err != nil {
		r.logger.Info("TypeCheck returns error", "err", err)
	}
	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*TimeNamingRule) Name() string {
	return "time-naming"
}

// SetLogger sets the logger field.
func (r *TimeNamingRule) SetLogger(logger *slog.Logger) {
	if logger != nil {
		r.logger = logger.With("rule", r.Name())
	}
}

type lintTimeNames struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintTimeNames) Visit(node ast.Node) ast.Visitor {
	v, ok := node.(*ast.ValueSpec)
	if !ok {
		return w
	}
	for _, name := range v.Names {
		origTyp := w.file.Pkg.TypeOf(name)
		// Look for time.Duration or *time.Duration;
		// the latter is common when using flag.Duration.
		typ := origTyp
		if pt, ok := typ.(*types.Pointer); ok {
			typ = pt.Elem()
		}
		if !isNamedType(typ, "time", "Duration") {
			continue
		}
		suffix := ""
		for _, suf := range timeSuffixes {
			if strings.HasSuffix(name.Name, suf) {
				suffix = suf
				break
			}
		}
		if suffix == "" {
			continue
		}
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryTime,
			Confidence: 0.9,
			Node:       v,
			Failure:    fmt.Sprintf("var %s is of type %v; don't use unit-specific suffix %q", name.Name, origTyp, suffix),
		})
	}
	return w
}

// timeSuffixes is a list of name suffixes that imply a time unit.
// This is not an exhaustive list.
var timeSuffixes = []string{
	"Hour", "Hours",
	"Min", "Mins", "Minutes", "Minute",
	"Sec", "Secs", "Seconds", "Second",
	"Msec", "Msecs",
	"Milli", "Millis", "Milliseconds", "Millisecond",
	"Usec", "Usecs", "Microseconds", "Microsecond",
	"MS", "Ms",
}

func isNamedType(typ types.Type, importPath, name string) bool {
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := n.Obj()
	return typeName != nil && typeName.Pkg() != nil && typeName.Pkg().Path() == importPath && typeName.Name() == name
}
