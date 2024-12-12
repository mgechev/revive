package rule

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// ImportsBlocklistRule disallows importing the specified packages.
type ImportsBlocklistRule struct {
	blocklist []*regexp.Regexp

	configureOnce sync.Once
	configureErr  error
}

var replaceImportRegexp = regexp.MustCompile(`/?\*\*/?`)

func (r *ImportsBlocklistRule) configure(arguments lint.Arguments) error {
	r.blocklist = []*regexp.Regexp{}
	for _, arg := range arguments {
		argStr, ok := arg.(string)
		if !ok {
			return fmt.Errorf("invalid argument to the imports-blocklist rule. Expecting a string, got %T", arg)
		}
		regStr, err := regexp.Compile(fmt.Sprintf(`(?m)"%s"$`, replaceImportRegexp.ReplaceAllString(argStr, `(\W|\w)*`)))
		if err != nil {
			return fmt.Errorf("invalid argument to the imports-blocklist rule. Expecting %q to be a valid regular expression, got: %w", argStr, err)
		}
		r.blocklist = append(r.blocklist, regStr)
	}
	return nil
}

func (r *ImportsBlocklistRule) isBlocklisted(path string) bool {
	for _, regex := range r.blocklist {
		if regex.MatchString(path) {
			return true
		}
	}
	return false
}

// Apply applies the rule to given file.
func (r *ImportsBlocklistRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	r.configureOnce.Do(func() { r.configureErr = r.configure(arguments) })
	if r.configureErr != nil {
		return newInternalFailureError(r.configureErr)
	}

	var failures []lint.Failure

	for _, is := range file.AST.Imports {
		path := is.Path
		if path != nil && r.isBlocklisted(path.Value) {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    "should not use the following blocklisted import: " + path.Value,
				Node:       is,
				Category:   "imports",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*ImportsBlocklistRule) Name() string {
	return "imports-blocklist"
}
