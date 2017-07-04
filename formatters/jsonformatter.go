package formatters

import "github.com/mgechev/golinter/visitors"
import "encoding/json"

type JSONFormatter struct {
	Metadata FormatterMetadata
}

// {
//   Name: "JSON Formatter",
//   Description: "This formatter produces JSON from the errors",
//   Sample: "[{ \"position\": 10, \"failure\": \"Forbidden semicolon\" }]"
// }

func (f *JSONFormatter) Format(failures []visitors.Failure) (string, error) {
	result, error := json.Marshal(failures)
	if error != nil {
		return "", error
	}
	return string(result), nil
}
