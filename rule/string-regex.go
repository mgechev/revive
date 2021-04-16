package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"github.com/mgechev/revive/lint"
)

// #region Revive API

// StringRegexRule lints strings and/or comments according to a set of regular expressions given as Arguments
type StringRegexRule struct{}

// Apply applies the rule to the given file.
func (r *StringRegexRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := lintStringRegexRule{onFailure: onFailure}
	w.parseArguments(arguments)
	ast.Walk(w, file.AST)

	return failures
}

func (r *StringRegexRule) Name() string {
	return "string-regex"
}

// #endregion

// #region Internal structure

type lintStringRegexRule struct {
	onFailure func(lint.Failure)

	rules              []stringRegexSubrule
	stringDeclarations map[string]string
}

type stringRegexSubrule struct {
	parent       *lintStringRegexRule
	scope        stringRegexSubruleScope
	regexp       *regexp.Regexp
	errorMessage string
}

type stringRegexSubruleScope struct {
	funcName string // Function name the rule is scoped to
	argument int    // (optional) Which argument in calls to the function is checked against the rule (the first argument is checked by default)
	field    string // (optional) If the argument to be checked is a struct, which member of the struct is checked against the rule (top level members only)
}

var parseStringRegexScope = regexp.MustCompile("^([A-Za-z][\\.A-Za-z0-9]+)(?:\\[([0-9]+)\\])?(?:\\.([A-Za-z]+))?$")

// #endregion

// #region Argument parsing

func (w *lintStringRegexRule) parseArguments(arguments lint.Arguments) {
	for i, argument := range arguments {
		scope, regex, errorMessage := w.parseArgument(argument, i)
		w.rules = append(w.rules, stringRegexSubrule{
			parent:       w,
			scope:        scope,
			regexp:       regex,
			errorMessage: errorMessage,
		})
	}
}

func (w lintStringRegexRule) parseArgument(argument interface{}, ruleNum int) (scope stringRegexSubruleScope, regex *regexp.Regexp, errorMessage string) {
	g, ok := argument.([]interface{}) // Cast to generic slice first
	if !ok {
		panic(fmt.Sprintf("unable to parse argument %d", ruleNum))
	}
	var rule []string
	for i, obj := range g {
		val, ok := obj.(string)
		if !ok {
			panic(fmt.Sprintf("unable to parse argument %d, option %d", ruleNum, i))
		}
		rule = append(rule, val)
	}

	// Parse rule scope
	scope = stringRegexSubruleScope{}
	matches := parseStringRegexScope.FindStringSubmatch(rule[0])
	if matches == nil {
		panic(fmt.Sprintf("unable to parse rule scope (argument %d, option 0)", ruleNum))
	}
	scope.funcName = matches[1]
	if len(matches[2]) > 0 {
		var err error
		scope.argument, err = strconv.Atoi(matches[2])
		if err != nil {
			fmt.Println(len(matches))
			panic(fmt.Sprintf("unable to parse rule scope argument number (argument %d, option 0)", ruleNum))
		}
	}
	if len(matches[3]) > 0 {
		scope.field = matches[3]
	}

	// Strip / characters from the beginning and end of rule[1] before compiling
	regex, err := regexp.Compile(rule[1][1 : len(rule[1])-1])
	if err != nil {
		panic(fmt.Sprintf("unable to compile %s as regexp (argument %d, option 1)", rule[1], ruleNum))
	}

	// Use custom error message if provided
	if len(rule) == 3 {
		errorMessage = rule[2]
	}
	return scope, regex, errorMessage
}

// #endregion

// #region Node traversal

func (w lintStringRegexRule) Visit(node ast.Node) ast.Visitor {
	// First, check if node is a call expression
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Get the name of the call expression to check against rule scope
	callName, ok := w.getCallName(call)
	if !ok {
		return w
	}

	for _, rule := range w.rules {
		if rule.scope.funcName == callName {
			rule.Apply(call)
		}
	}

	return w
}

// Return the name of a call expression in the form of package.Func or Func
func (w lintStringRegexRule) getCallName(call *ast.CallExpr) (callName string, ok bool) {
	if ident, ok := call.Fun.(*ast.Ident); ok {
		// Local function call
		return ident.Name, true
	} else if selector, ok := call.Fun.(*ast.SelectorExpr); ok {
		// Scoped function call
		scope, ok := selector.X.(*ast.Ident)
		if !ok {
			return "", false
		}
		return scope.Name + "." + selector.Sel.Name, true
	} else {
		return "", false
	}
}

// #endregion

// #region Linting logic

// Apply a single regex rule to a call expression (should be done after verifying the that the call expression matches the rule's scope)
func (rule stringRegexSubrule) Apply(call *ast.CallExpr) {
	if len(call.Args) <= rule.scope.argument {
		// TODO: should cases where calls are incompatible with scope cause failures?
		return
	}

	arg := call.Args[rule.scope.argument]
	var lit *ast.BasicLit
	if len(rule.scope.field) > 0 {
		// Try finding the scope's Field, treating arg as a composite literal
		composite, ok := arg.(*ast.CompositeLit)
		if !ok {
			return
		}
		for _, el := range composite.Elts {
			kv, ok := el.(*ast.KeyValueExpr)
			if !ok {
				continue
			}
			key, ok := kv.Key.(*ast.Ident)
			if !ok || key.Name != rule.scope.field {
				continue
			}

			// We're now dealing with the exact field in the rule's scope, so if anything fails, we can safely return instead of continuing the loop
			lit, ok = kv.Value.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return
			}
		}
	} else {
		var ok bool
		// Treat arg as a string literal
		lit, ok = arg.(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return
		}
	}
	// Unquote the string literal before linting
	unquoted := lit.Value[1 : len(lit.Value)-1]
	rule.lintMessage(unquoted, lit)
}

func (rule stringRegexSubrule) lintMessage(s string, node ast.Node) {
	// Fail if the string doesn't match the user's regex
	if rule.regexp.MatchString(s) {
		return
	}
	var failure string
	if len(rule.errorMessage) > 0 {
		failure = fmt.Sprintf("string literal doesn't match user defined regex (%s)", rule.errorMessage)
	} else {
		failure = fmt.Sprintf("string literal doesn't match user defined regex /%s/", rule.regexp.String())
	}
	rule.parent.onFailure(lint.Failure{
		Confidence: 1,
		Failure:    failure,
		Node:       node})
}

// #endregion
