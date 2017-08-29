package formatter

import (
	"encoding/json"

	"github.com/mgechev/revive/rule"
)

// JSONFormatter is an implementation of the Formatter interface
// which formats the errors to JSON.
type JSONFormatter struct {
	Metadata FormatterMetadata
}

// Format formats the failures gotten from the linter.
func (f *JSONFormatter) Format(failures []rule.Failure) (string, error) {
	result, error := json.Marshal(failures)
	if error != nil {
		return "", error
	}
	return string(result), nil
}
