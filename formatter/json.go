package formatter

import (
	"encoding/json"

	"github.com/mgechev/revive/lint"
)

// JSON is an implementation of the Formatter interface
// which formats the errors to JSON.
type JSON struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter
func (f *JSON) Name() string {
	return "json"
}

// Format formats the failures gotten from the lint.
func (f *JSON) Format(failures <-chan lint.Failure, config lint.RulesConfig) (string, error) {
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
