package formatter

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mgechev/revive/lint"
)

// Stylish is an implementation of the Formatter interface
// which formats the errors to JSON.
type Stylish struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter
func (*Stylish) Name() string {
	return "stylish"
}

func formatFailure(failure lint.Failure, severity lint.Severity) []string {
	fString := color.CyanString(failure.Failure)
	fURL := ruleDescriptionURL(failure.RuleName)
	fName := color.RedString(fURL)
	lineColumn := failure.Position
	pos := fmt.Sprintf("(%d, %d)", lineColumn.Start.Line, lineColumn.Start.Column)
	if severity == lint.SeverityWarning {
		fName = color.YellowString(fURL)
	}
	return []string{failure.GetFilename(), pos, fName, fString}
}

// Format formats the failures gotten from the lint.
func (s *Stylish) Format(failures <-chan lint.Failure, config lint.Config) (string, error) {
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
	ps := "problems"
	if total == 1 {
		ps = "problem"
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
		output += s.table(val) + "\n"
	}

	suffix := fmt.Sprintf(" %d %s (%d errors) (%d warnings)", total, ps, totalErrors, total-totalErrors)

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

func (*Stylish) table(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	colWidths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, col := range row {
			if w := utf8.RuneCountInString(col); w > colWidths[i] {
				colWidths[i] = w
			}
		}
	}

	var buf strings.Builder
	for _, row := range rows {
		buf.WriteString("  ")
		for i, col := range row {
			fmt.Fprintf(&buf, "%-*s", colWidths[i]+2, col)
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}
