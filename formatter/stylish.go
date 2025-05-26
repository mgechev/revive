package formatter

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/mgechev/revive/lint"
)

// Stylish is an implementation of the Formatter interface
// which formats the errors to JSON.
type Stylish struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter.
func (*Stylish) Name() string {
	return "stylish"
}

func formatFailure(failure lint.Failure, severity lint.Severity) []string {
	fString := color.CyanString(failure.Failure)
	lineColumn := failure.Position
	pos := fmt.Sprintf("(%d, %d)", lineColumn.Start.Line, lineColumn.Start.Column)
	fURL := ruleDescriptionURL(failure.RuleName)
	fName := color.RedString(fURL)
	if severity == lint.SeverityWarning {
		fName = color.YellowString(fURL)
	}
	return []string{failure.Filename(), pos, fName, fString}
}

// Format formats the failures gotten from the lint.
func (*Stylish) Format(failures <-chan lint.Failure, config lint.Config) (string, error) {
	var result [][]string
	totalErrors := 0
	total := 0

	for f := range failures {
		total++
		currentType := severity(config, f)
		if currentType == lint.SeverityError {
			totalErrors++
		}
		result = append(result, formatFailure(f, lint.Severity(currentType)))
	}

	fileReport := map[string][][]string{}

	for _, row := range result {
		if _, ok := fileReport[row[0]]; !ok {
			fileReport[row[0]] = [][]string{}
		}

		fileReport[row[0]] = append(fileReport[row[0]], []string{row[1], row[2], row[3]})
	}

	output := ""
	for filename, val := range fileReport {
		c := color.New(color.Underline)
		output += c.SprintfFunc()(filename + "\n")
		output += table(val) + "\n"
	}

	problemsLabel := "problems"
	if total == 1 {
		problemsLabel = "problem"
	}
	totalWarnings := total - totalErrors
	warningsLabel := "warnings"
	if totalWarnings == 1 {
		warningsLabel = "warning"
	}
	errorsLabel := "errors"
	if totalErrors == 1 {
		errorsLabel = "error"
	}
	suffix := fmt.Sprintf(" %d %s (%d %s) (%d %s)", total, problemsLabel, totalErrors, errorsLabel, totalWarnings, warningsLabel)

	switch {
	case total > 0 && totalErrors > 0:
		suffix = color.RedString("\n ✖" + suffix)
	case total > 0 && totalErrors == 0:
		suffix = color.YellowString("\n ✖" + suffix)
	default:
		suffix, output = "", ""
	}

	return output + suffix, nil
}
