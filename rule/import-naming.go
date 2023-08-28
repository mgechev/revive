package rule

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// ImportNamingRule lints import naming.
type ImportNamingRule struct {
	configured       bool
	namingRule       string
	namingRuleRegexp *regexp.Regexp
	sync.Mutex
}

const defaultNamingRule = "^[a-z][a-z0-9]$"

var defaultNamingRuleRegexp = regexp.MustCompile(defaultNamingRule)

func (r *ImportNamingRule) configure(arguments lint.Arguments) {
	r.Lock()
	defer r.Unlock()
	if r.configured {
		return
	}

	if len(arguments) < 1 {
		r.namingRule = defaultNamingRule
		r.namingRuleRegexp = defaultNamingRuleRegexp
		return
	}

	var ok bool
	r.namingRule, ok = arguments[0].(string) // Alt. non panicking version
	if !ok {
		panic(`invalid value passed as argument number to the "import-naming" rule`)
	}

	var err error
	r.namingRuleRegexp, err = regexp.Compile(r.namingRule)
	if err != nil {
		panic(fmt.Sprintf("Invalid argument to the import-naming rule. Expecting %q to be a valid regular expression, got: %v", r.namingRule, err))
	}
}

// Apply applies the rule to given file.
func (r *ImportNamingRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configure(arguments)

	var failures []lint.Failure

	for _, is := range file.AST.Imports {
		path := is.Path
		if path == nil {
			continue
		}

		alias := is.Name
		if alias == nil || alias.Name == "_" {
			continue
		}

		if !r.namingRuleRegexp.MatchString(alias.Name) {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("import name (%s) must match the regular expression: %s", alias.Name, r.namingRule),
				Node:       alias,
				Category:   "imports",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*ImportNamingRule) Name() string {
	return "import-naming"
}
