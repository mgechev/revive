package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"github.com/mgechev/revive/lint"
)

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

type lintStringRegexRule struct {
	onFailure func(lint.Failure)

	rules []stringRegexSubrule
}

type stringRegexSubrule struct {
	Scope        stringRegexScope
	Regexp       *regexp.Regexp
	ErrorMessage *string
}

type stringRegexScope struct {
	Func     string // Function name the rule is scoped to
	Argument int    // (optional) Which argument in calls to the function is checked against the rule (the first argument is checked by default)
	Member   string // (optional) If the argument to be checked is a struct, which member of the struct is checked against the rule (top level members only)
}

var checkStringRegexScope = struct {
	Basic              *regexp.Regexp
	WithArgument       *regexp.Regexp
	WithArgumentMember *regexp.Regexp
}{
	Basic:              regexp.MustCompile("^[A-Za-z]+"),
	WithArgument:       regexp.MustCompile("\\[\\d\\]"),
	WithArgumentMember: regexp.MustCompilePOSIX("^\\[\\d\\]\\.\\K[A-Za-z]+$"),
}

func (w lintStringRegexRule) Visit(node ast.Node) ast.Visitor {
	// First, check if node is a string literal
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	// Unquote the string
	unquoted := lit.Value[1 : len(lit.Value)-1]
	w.lintMessage(unquoted, node)

	return w
}

func (w *lintStringRegexRule) parseArguments(arguments lint.Arguments) {
	for i, argument := range arguments {
		scope, regex, errorMessage := w.parseArgument(argument, i)
		w.rules = append(w.rules, stringRegexSubrule{
			Scope:        scope,
			Regexp:       regex,
			ErrorMessage: errorMessage,
		})
	}
}

func (w lintStringRegexRule) parseArgument(argument interface{}, ruleNum int) (scope stringRegexScope, regex *regexp.Regexp, errorMessage *string) {
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
	scope = stringRegexScope{}
	c := checkStringRegexScope
	if funcName := c.Basic.FindString(rule[0]); len(funcName) > 0 {
		scope.Func = funcName
	} else {
		panic(fmt.Sprintf("rule scope doesn't start with a valid function name (argument %d, option 0)", ruleNum))
	}
	if arg := c.WithArgument.FindString(rule[0]); len(arg) > 0 {
		argNum, err := strconv.Atoi(arg)
		if err != nil {
			panic(fmt.Sprintf("invalid argument number given in rule scope (argument %d, option 0)", ruleNum))
		}
		scope.Argument = argNum
	}
	if member := c.WithArgumentMember.FindString(rule[0]); len(member) > 0 {
		scope.Member = member
	}

	// Strip / characters from the beginning and end of rule[1] before compiling
	regex, err := regexp.Compile(rule[1][1 : len(rule[1])-1])
	if err != nil {
		panic(fmt.Sprintf("unable to compile %s as regexp (argument %d, option 1)", rule[1], ruleNum))
	}

	// Parse custom error message if provided
	if len(rule) == 2 {
		errorMessage = &rule[1]
	}
	return scope, regex, errorMessage
}

func (w lintStringRegexRule) lintMessage(s string, node ast.Node) {
	for _, rule := range w.rules {
		// Fail if the string doesn't match the user's regex
		if rule.Regexp.Match([]byte(s)) {
			continue
		}
		var failure string
		if rule.ErrorMessage != nil {
			failure = fmt.Sprintf("string literal doesn't match user defined regex (%s)", *rule.ErrorMessage)
		} else {
			failure = fmt.Sprintf("string literal doesn't match user defined regex /%s/", rule.Regexp.String())
		}
		w.onFailure(lint.Failure{
			Confidence: 1,
			Failure:    failure,
			Node:       node})
	}
}
