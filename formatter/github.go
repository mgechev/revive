package formatter

import (
	"bytes"
	"fmt"

	"github.com/mgechev/revive/lint"
)

// Github is an implementation of the Formatter interface
// which formats the errors for GitHub Actions annotations.
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message
type Github struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter
func (*Github) Name() string {
	return "github"
}

// Format formats the failures gotten from the lint.
func (*Github) Format(failures <-chan lint.Failure, config lint.Config) (string, error) {
	var buf bytes.Buffer
	for failure := range failures {
		fmt.Fprintf(&buf, "::%s file=%s,line=%d,col=%d,endLine=%d,endCol=%d,title=Revive: %s::%s",
			severity(config, failure),
			failure.Position.Start.Filename,
			failure.Position.Start.Line,
			failure.Position.Start.Column,
			failure.Position.End.Line,
			failure.Position.End.Column,
			failure.RuleName,
			failure.Failure,
		)
	}

	return buf.String(), nil
}
