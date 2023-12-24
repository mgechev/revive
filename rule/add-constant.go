package rule

import (
	"fmt"
	"go/ast"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

const (
	defaultStrLitLimit = 2
	kindFLOAT          = "FLOAT"
	kindINT            = "INT"
	kindSTRING         = "STRING"
)

type whiteList map[string]map[string]bool

func newWhiteList() whiteList {
	return map[string]map[string]bool{kindINT: {}, kindFLOAT: {}, kindSTRING: {}}
}

func (wl whiteList) add(kind, list string) {
	elems := strings.Split(list, ",")
	for _, e := range elems {
		wl[kind][e] = true
	}
}

// AddConstantRule lints unused params in functions.
type AddConstantRule struct {
	whiteList       whiteList
	ignoreFunctions []*regexp.Regexp
	strLitLimit     int
	sync.Mutex
}

// Apply applies the rule to given file.
func (r *AddConstantRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configure(arguments)

	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintAddConstantRule{
		onFailure:       onFailure,
		strLits:         make(map[string]int),
		strLitLimit:     r.strLitLimit,
		whiteLst:        r.whiteList,
		ignoreFunctions: r.ignoreFunctions,
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*AddConstantRule) Name() string {
	return "add-constant"
}

type lintAddConstantRule struct {
	onFailure       func(lint.Failure)
	strLits         map[string]int
	strLitLimit     int
	whiteLst        whiteList
	ignoreFunctions []*regexp.Regexp
}

func (r lintAddConstantRule) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.CallExpr:
		r.checkFunc(n)
		return nil
	case *ast.GenDecl:
		return nil // skip declarations
	case *ast.BasicLit:
		r.checkLit(n)
	}

	return r
}

func (r lintAddConstantRule) checkFunc(expr *ast.CallExpr) {
	fName := r.getFuncName(expr)

	for _, arg := range expr.Args {
		switch t := arg.(type) {
		case *ast.CallExpr:
			r.checkFunc(t)
		case *ast.BasicLit:
			if r.isIgnoredFunc(fName) {
				continue
			}
			r.checkLit(t)
		}
	}
}

func (lintAddConstantRule) getFuncName(expr *ast.CallExpr) string {
	switch f := expr.Fun.(type) {
	case *ast.SelectorExpr:
		switch prefix := f.X.(type) {
		case *ast.Ident:
			return prefix.Name + "." + f.Sel.Name
		}
	case *ast.Ident:
		return f.Name
	}

	return ""
}

func (r lintAddConstantRule) checkLit(n *ast.BasicLit) {
	switch kind := n.Kind.String(); kind {
	case kindFLOAT, kindINT:
		r.checkNumLit(kind, n)
	case kindSTRING:
		r.checkStrLit(n)
	}
}

func (r lintAddConstantRule) isIgnoredFunc(fName string) bool {
	for _, pattern := range r.ignoreFunctions {
		if pattern.MatchString(fName) {
			return true
		}
	}

	return false
}

func (r lintAddConstantRule) checkStrLit(n *ast.BasicLit) {
	if r.whiteLst[kindSTRING][n.Value] {
		return
	}

	count := r.strLits[n.Value]
	if count >= 0 {
		r.strLits[n.Value] = count + 1
		if r.strLits[n.Value] > r.strLitLimit {
			r.onFailure(lint.Failure{
				Confidence: 1,
				Node:       n,
				Category:   "style",
				Failure:    fmt.Sprintf("string literal %s appears, at least, %d times, create a named constant for it", n.Value, r.strLits[n.Value]),
			})
			r.strLits[n.Value] = -1 // mark it to avoid failing again on the same literal
		}
	}
}

func (r lintAddConstantRule) checkNumLit(kind string, n *ast.BasicLit) {
	if r.whiteLst[kind][n.Value] {
		return
	}

	r.onFailure(lint.Failure{
		Confidence: 1,
		Node:       n,
		Category:   "style",
		Failure:    fmt.Sprintf("avoid magic numbers like '%s', create a named constant for it", n.Value),
	})
}

func (r *AddConstantRule) configure(arguments lint.Arguments) {
	r.Lock()
	defer r.Unlock()

	if r.whiteList == nil {
		r.strLitLimit = defaultStrLitLimit
		r.whiteList = newWhiteList()
		if len(arguments) > 0 {
			args, ok := arguments[0].(map[string]any)
			if !ok {
				panic(fmt.Sprintf("Invalid argument to the add-constant rule. Expecting a k,v map, got %T", arguments[0]))
			}
			for k, v := range args {
				kind := ""
				switch k {
				case "allowFloats":
					kind = kindFLOAT
					fallthrough
				case "allowInts":
					if kind == "" {
						kind = kindINT
					}
					fallthrough
				case "allowStrs":
					if kind == "" {
						kind = kindSTRING
					}
					list, ok := v.(string)
					if !ok {
						panic(fmt.Sprintf("Invalid argument to the add-constant rule, string expected. Got '%v' (%T)", v, v))
					}
					r.whiteList.add(kind, list)
				case "maxLitCount":
					sl, ok := v.(string)
					if !ok {
						panic(fmt.Sprintf("Invalid argument to the add-constant rule, expecting string representation of an integer. Got '%v' (%T)", v, v))
					}

					limit, err := strconv.Atoi(sl)
					if err != nil {
						panic(fmt.Sprintf("Invalid argument to the add-constant rule, expecting string representation of an integer. Got '%v'", v))
					}
					r.strLitLimit = limit
				case "ignoreFuncs":
					excludes, ok := v.(string)
					if !ok {
						panic(fmt.Sprintf("Invalid argument to the ignoreFuncs parameter of add-constant rule, string expected. Got '%v' (%T)", v, v))
					}

					for _, exclude := range strings.Split(excludes, ",") {
						exclude = strings.Trim(exclude, " ")
						if exclude == "" {
							panic("Invalid argument to the ignoreFuncs parameter of add-constant rule, expected regular expression must not be empty.")
						}

						exp, err := regexp.Compile(exclude)
						if err != nil {
							panic(fmt.Sprintf("Invalid argument to the ignoreFuncs parameter of add-constant rule: regexp %q does not compile: %v", exclude, err))
						}

						r.ignoreFunctions = append(r.ignoreFunctions, exp)
					}
				}
			}
		}
	}
}
