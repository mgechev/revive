package formatters

import "github.com/mgechev/golinter/visitors"

type FormatterMetadata struct {
	Name        string
	Description string
	Sample      string
}

type Formatter interface {
	Format([]visitors.Failure) string
}
