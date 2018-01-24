package formatter

import (
	"bytes"
	"fmt"

	"github.com/fatih/color"
	"github.com/mgechev/revive/linter"
	"github.com/olekukonko/tablewriter"
)

const (
	errorEmoji   = ""
	warningEmoji = ""
)

// CLIFormatter is an implementation of the Formatter interface
// which formats the errors to JSON.
type CLIFormatter struct {
	Metadata linter.FormatterMetadata
}

func formatFailure(failure linter.Failure, severity linter.Severity) []string {
	fString := color.BlueString(failure.Failure)
	fTypeStr := string(severity)
	fType := color.RedString(fTypeStr)
	lineColumn := failure.Position
	pos := fmt.Sprintf("(%d, %d)", lineColumn.Start.Line, lineColumn.Start.Column)
	if severity == linter.SeverityWarning {
		fType = color.YellowString(fTypeStr)
	}
	return []string{failure.GetFilename(), pos, fType, fString}
}

// Format formats the failures gotten from the linter.
func (f *CLIFormatter) Format(failures <-chan linter.Failure, config linter.RulesConfig) (string, error) {
	var result [][]string
	var totalErrors = 0
	var total = 0

	for f := range failures {
		total++
		currentType := linter.SeverityWarning
		if config, ok := config[f.RuleName]; ok && config.Severity == linter.SeverityError {
			currentType = linter.SeverityError
			totalErrors++
		}
		result = append(result, formatFailure(f, linter.Severity(currentType)))
	}
	ps := "problems"
	if total == 1 {
		ps = "problem"
	}

	fileReport := make(map[string][][]string)

	for _, row := range result {
		if _, ok := fileReport[row[0]]; !ok {
			fileReport[row[0]] = [][]string{}
		}

		fileReport[row[0]] = append(fileReport[row[0]], []string{row[1], row[2], row[3]})
	}

	output := ""
	for filename, val := range fileReport {
		buf := new(bytes.Buffer)
		table := tablewriter.NewWriter(buf)
		table.SetBorder(false)
		table.SetColumnSeparator("")
		table.SetRowSeparator("")
		table.SetAutoWrapText(false)
		table.AppendBulk(val)
		table.Render()
		c := color.New(color.Underline)
		output += c.SprintfFunc()(filename + "\n")
		output += buf.String() + "\n"
	}

	suffix := fmt.Sprintf(" %d %s (%d errors) (%d warnings)", total, ps, totalErrors, total-totalErrors)

	if total > 0 && totalErrors > 0 {
		suffix = color.RedString("\n ✖" + suffix)
	} else if total > 0 && totalErrors == 0 {
		suffix = color.YellowString("\n ✖" + suffix)
	} else {
		suffix = color.GreenString("\n" + suffix)
	}
	return output + suffix, nil
}
