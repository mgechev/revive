package formatter

import (
	"encoding/json"
	"os"

	"github.com/mgechev/revive/lint"
)

// JSONStream is an implementation of the Formatter interface
// which formats the errors to JSON.
type JSONStream struct {
	Metadata lint.FormatterMetadata
}

// Name returns the name of the formatter
func (f *JSONStream) Name() string {
	return "json-stream"
}

// Format formats the failures gotten from the lint.
func (f *JSONStream) Format(failures <-chan lint.Failure, _ lint.RulesConfig) (string, error) {
	enc := json.NewEncoder(os.Stdout)
	for failure := range failures {
		err := enc.Encode(failure)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}
