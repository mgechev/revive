package rule

import (
	"fmt"
	"sync"

	"github.com/mgechev/revive/lint"
)

// ImportsBlacklistRule lints given else constructs.
type ImportsBlacklistRule struct {
	blacklist map[string]bool
	sync.Mutex
}

func (r *ImportsBlacklistRule) configure(arguments lint.Arguments) {
	r.Lock()
	if r.blacklist == nil {
		r.blacklist = make(map[string]bool, len(arguments))

		for _, arg := range arguments {
			argStr, ok := arg.(string)
			if !ok {
				panic(fmt.Sprintf("Invalid argument to the imports-blacklist rule. Expecting a string, got %T", arg))
			}
			// we add quotes if not present, because when parsed, the value of the AST node, will be quoted
			if len(argStr) > 2 && argStr[0] != '"' && argStr[len(argStr)-1] != '"' {
				argStr = fmt.Sprintf(`%q`, argStr)
			}
			r.blacklist[argStr] = true
		}
	}
	r.Unlock()
}

// Apply applies the rule to given file.
func (r *ImportsBlacklistRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configure(arguments)

	var failures []lint.Failure

	if file.IsTest() {
		return failures // skip, test file
	}

	for _, is := range file.AST.Imports {
		path := is.Path
		if path != nil && r.blacklist[path.Value] {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    "should not use the following blacklisted import: " + path.Value,
				Node:       is,
				Category:   "imports",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (r *ImportsBlacklistRule) Name() string {
	return "imports-blacklist"
}
