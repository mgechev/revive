package linter

import (
	"go/token"

	"github.com/mgechev/revive/file"
	"github.com/mgechev/revive/rule"
)

// ReadFile defines an abstraction for reading files.
type ReadFile func(path string) (result []byte, err error)

// Linter is used for lintign set of files.
type Linter struct {
	reader ReadFile
}

// New creates a new Linter
func New(reader ReadFile) Linter {
	return Linter{reader: reader}
}

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(filenames []string, ruleSet []rule.Rule) ([]rule.Failure, error) {
	var fileSet token.FileSet
	var failures []rule.Failure
	for _, filename := range filenames {
		content, err := l.reader(filename)
		if err != nil {
			return nil, err
		}
		file, err := file.New(filename, content, &fileSet)

		if err != nil {
			return nil, err
		}

		for _, rule := range ruleSet {
			currentFailures := rule.Apply(file, []string{})
			failures = append(failures, currentFailures...)
		}
	}

	return failures, nil
}
