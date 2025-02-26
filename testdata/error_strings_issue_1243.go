package fixtures

import (
	"errors"
	"fmt"
)

func issue1243() {
	err := errors.New("An error occurred!") // MATCH /error strings should not be capitalized or end with punctuation or a newline/
	fmt.Println(err)
}
