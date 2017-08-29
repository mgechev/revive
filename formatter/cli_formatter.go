package formatter

import (
	"bytes"
	"fmt"

	"github.com/mgechev/revive/rule"
	"github.com/olekukonko/tablewriter"
	"github.com/ttacon/chalk"
)

const (
	errorEmoji   = ""
	warningEmoji = ""
)

// CLIFormatter is an implementation of the Formatter interface
// which formats the errors to JSON.
type CLIFormatter struct {
	Metadata FormatterMetadata
}

func formatFailure(failure rule.Failure) []string {
	fString := chalk.Blue.Color(failure.Failure)
	fTypeStr := string(failure.Type)
	fType := chalk.Red.Color(fTypeStr)
	lineColumn := failure.Position
	pos := chalk.Dim.TextStyle(fmt.Sprintf("(%d, %d)", lineColumn.Start.Line, lineColumn.Start.Column))
	if failure.Type == rule.FailureTypeWarning {
		fType = chalk.Yellow.Color(fTypeStr)
	}
	return []string{failure.GetFilename(), pos, fType, fString}
}

// Format formats the failures gotten from the linter.
func (f *CLIFormatter) Format(failures []rule.Failure) (string, error) {
	var result [][]string
	var totalErrors = 0
	for _, f := range failures {
		result = append(result, formatFailure(f))
		if f.Type == rule.FailureTypeError {
			totalErrors++
		}
	}
	total := len(failures)
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
		output += chalk.Dim.TextStyle(chalk.Underline.TextStyle(filename) + "\n")
		output += buf.String() + "\n"
	}

	suffix := fmt.Sprintf(" %d %s (%d errors) (%d warnings)", total, ps, totalErrors, total-totalErrors)

	if total > 0 && totalErrors > 0 {
		suffix = chalk.Red.Color("\n ✖" + suffix)
	} else if total > 0 && totalErrors == 0 {
		suffix = chalk.Yellow.Color("\n ✖" + suffix)
	} else {
		suffix = chalk.Green.Color("\n" + suffix)
	}
	return output + suffix, nil
}
