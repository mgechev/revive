package formatters

import "github.com/mgechev/golinter/rules"

// FormatterMetadata configuration of a formatter
type FormatterMetadata struct {
	Name        string
	Description string
	Sample      string
}

// Formatter defines an interface for failure formatters
type Formatter interface {
	Format([]rules.Failure) string
}
