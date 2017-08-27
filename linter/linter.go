package linter

import (
	"go/token"

	"github.com/mgechev/golinter/file"
	"github.com/mgechev/golinter/rules"
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

// Lint lints a set of files with the specified rules.
func (l *Linter) Lint(filenames []string, ruleSet []rules.Rule) ([]rules.Failure, error) {
	var fileSet token.FileSet
	var failures []rules.Failure
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
