package formatter

import (
	"encoding/json"

	"github.com/mgechev/revive/lint"
)

// JSONFormatter is an implementation of the Formatter interface
// which formats the errors to JSON.
type JSONFormatter struct {
	Metadata lint.FormatterMetadata
}

// Format formats the failures gotten from the lint.
func (f *JSONFormatter) Format(failures <-chan lint.Failure, config lint.RulesConfig) (string, error) {
	var slice []lint.Failure
	for failure := range failures {
		slice = append(slice, failure)
	}
	result, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(result), err
}
