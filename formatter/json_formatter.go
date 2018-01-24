package formatter

import (
	"encoding/json"

	"github.com/mgechev/revive/linter"
)

// JSONFormatter is an implementation of the Formatter interface
// which formats the errors to JSON.
type JSONFormatter struct {
	Metadata linter.FormatterMetadata
}

// Format formats the failures gotten from the linter.
func (f *JSONFormatter) Format(failures <-chan linter.Failure, config linter.RulesConfig) (string, error) {
	var slice []linter.Failure
	for failure := range failures {
		slice = append(slice, failure)
	}
	result, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(result), err
}
