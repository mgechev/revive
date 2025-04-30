package formatter

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mgechev/revive/lint"
)

// GitHub is an implementation of the Formatter interface
// which formats the errors for GitHub Actions annotations.
// https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions#setting-an-error-message
type GitHub struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter
func (*GitHub) Name() string {
	return "github"
}

// Format formats the failures gotten from the lint.
func (*GitHub) Format(failures <-chan lint.Failure, config lint.Config) (string, error) {
	replacer := strings.NewReplacer(
		"\n", "%0A",
		"\r", "%0D",
	)
	var buf bytes.Buffer
	for failure := range failures {
		message := replacer.Replace(failure.Failure)
		fmt.Fprintf(&buf, "::%s file=%s,line=%d,col=%d,endLine=%d,endCol=%d,title=Revive: %s::%s\n",
			severity(config, failure),
			failure.Position.Start.Filename,
			failure.Position.Start.Line,
			failure.Position.Start.Column,
			failure.Position.End.Line,
			failure.Position.End.Column,
			failure.RuleName,
			message,
		)
	}

	return buf.String(), nil
}
