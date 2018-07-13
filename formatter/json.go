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

// JSONObject defines a JSON object of an failure
type JSONObject struct {
	Severity     lint.Severity
	lint.Failure `json:",inline"`
}

// Format formats the failures gotten from the lint.
func (f *JSON) Format(failures <-chan lint.Failure, config lint.RulesConfig) (string, error) {
	var slice []JSONObject
	for failure := range failures {
		obj := JSONObject{}
		obj.Severity = severity(config, failure)
		obj.Failure = failure
		slice = append(slice, obj)
	}
	result, err := json.Marshal(slice)
	if err != nil {
		return "", err
	}
	return string(result), err
}
