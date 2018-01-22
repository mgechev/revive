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
func (f *JSONFormatter) Format(failures []linter.Failure) (string, error) {
	result, error := json.Marshal(failures)
	if error != nil {
		return "", error
	}
	return string(result), nil
}
