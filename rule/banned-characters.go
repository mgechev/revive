package rule

import (
	"bytes"
	"fmt"

	"github.com/mgechev/revive/lint"
)

// BannedCharsRule checks if a file contains banned characters.
type BannedCharsRule struct{}

const bannedCharsRuleName = "banned-characters"

// Apply applied the rule to the given file.
func (r *BannedCharsRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	checkNumberOfArguments(1, arguments, bannedCharsRuleName)
	bannedCharList := r.getBannedList(arguments)

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	// check if file content holds banned characters
	for _, c := range bannedCharList {
		ok := bytes.ContainsAny(file.Content(), c)
		if ok {
			onFailure(lint.Failure{
				Failure: fmt.Sprintf("banned character found: %s", c),
			})
		}
	}

	return failures
}

// Name returns the rule name
func (r *BannedCharsRule) Name() string {
	return bannedCharsRuleName
}

// getBannedList converts arguments into the banned characters list
func (r *BannedCharsRule) getBannedList(args lint.Arguments) []string {
	var bannedChars []string
	for _, char := range args {
		charStr, ok := char.(string)
		if !ok {
			panic(fmt.Sprintf("Invalid argument for the %s rule: expecting a string, got %T", r.Name(), char))
		}
		bannedChars = append(bannedChars, charStr)
	}

	return bannedChars
}
