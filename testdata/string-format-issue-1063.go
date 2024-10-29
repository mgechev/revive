package fixtures

import (
	"errors"
)

func ReturnError() error {
	return errors.New("This is an error.") // MATCH /must not start with a capital letter/
}
