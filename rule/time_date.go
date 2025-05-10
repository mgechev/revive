package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"log/slog"
	"strconv"
	"strings"

	"github.com/mgechev/revive/lint"
)

// TimeDateRule lints the way time.Date is used.
type TimeDateRule struct{}

// Apply applies the rule to given file.
func (*TimeDateRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTimeDate{file, onFailure}

	ast.Walk(w, file.AST)
	return failures
}

// Name returns the rule name.
func (*TimeDateRule) Name() string {
	return "time-date"
}

type lintTimeDate struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

var (
	errTimeDateFoundOctal            = errors.New("octal notation found")
	errTimeDateFoundOctalLeadingZero = errors.New("octal notation with leading zero found")
	errTimeDateFoundOctalPaddingZero = errors.New("octal notation with padding zeroes found")
	errTimeDateFoundHexadecimal      = errors.New("hexadecimal notation found")
	errTimeDateFoundBinary           = errors.New("binary notation found")
	errTimeDateFoundFloat            = errors.New("float literal found")
	errTimeDateFoundExponential      = errors.New("exponential notation found")
	errTimeDateFoundAlternative      = errors.New("alternative notation found")
)

func (w lintTimeDate) Visit(n ast.Node) ast.Visitor {
	const timeDateArity = 8
	ce, ok := n.(*ast.CallExpr)
	if !ok || len(ce.Args) != timeDateArity {
		return w
	}
	if !isPkgDot(ce.Fun, "time", "Date") {
		return w
	}

	argumentsNames := []string{
		"year",
		"month",
		"day",
		"hour",
		"minute",
		"second",
		"nanosecond",
		"timezone",
	}

	// The last argument is a timezone, the is no need to check it
	for pos, arg := range ce.Args[:timeDateArity-1] {
		bl, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}

		replacedValue, err := w.validateInteger(bl)
		if err == nil {
			continue
		}

		confidence := 0.8 // default confidence
		if errors.Is(err, errTimeDateFoundOctalLeadingZero) {
			confidence = 0.5 // people can use 00, 01, 02, 03, 04, 05, 06, 07 if they want
		} else if errors.Is(err, errTimeDateFoundOctalPaddingZero) {
			confidence = 1 // this is a clear mistake
		}

		w.onFailure(lint.Failure{
			Category:   "time",
			Node:       bl,
			Confidence: confidence,
			Failure: fmt.Sprintf(
				"use decimal digits for time.Date %s argument: %s: use %s instead of %s",
				argumentsNames[pos], err, replacedValue, bl.Value),
		})
	}

	return w
}

func (lintTimeDate) validateInteger(bl *ast.BasicLit) (string, error) {
	currentValue := strings.ToLower(bl.Value)

	// handle the obvious case
	switch currentValue {
	case "0":
		// 0 is a valid value for all the arguments
		// let's skip it
		return bl.Value, nil
	case "00", "01", "02", "03", "04", "05", "06", "07":
		// someone used a leading zero, while they should have used a integer literal
		cleanedValue := currentValue[1:]
		return cleanedValue, errTimeDateFoundOctalLeadingZero
	case "0o0", "0o1", "0o2", "0o3", "0o4", "0o5", "0o6", "0o7":
		// octal literal notation can be found after using gofumpt
		cleanedValue := currentValue[2:]
		return cleanedValue, errTimeDateFoundOctal
	}

	switch bl.Kind {
	case token.FLOAT:
		// someone used a float literal, while they should have used a integer literal
		parsedValue, err := strconv.ParseFloat(currentValue, 64)
		if err != nil {
			slog.Debug("failed to parse number", slog.String("value", currentValue), slog.String("error", err.Error()))
			return bl.Value, nil
		}

		// this will convert back the number to a string
		cleanedValue := strconv.FormatFloat(parsedValue, 'f', -1, 64)
		if strings.Contains(currentValue, "e") {
			return cleanedValue, errTimeDateFoundExponential
		}

		return cleanedValue, errTimeDateFoundFloat

	case token.INT:
		// we expect this format

	default:
		// anything else is invalid
		slog.Debug("unexpected kind of literal", slog.String("value", currentValue), "kind", bl.Kind.String())

		// let's be defensive and return nil
		return bl.Value, nil
	}

	// Parse the number with base=0 that allows to accept all number formats and base
	parsedValue, err := strconv.ParseInt(currentValue, 0, 64)
	if err != nil {
		slog.Debug("failed to parse number", slog.String("value", currentValue), slog.String("error", err.Error()))
		return bl.Value, nil
	}

	cleanedValue := strconv.FormatInt(parsedValue, 10)

	// The number is not in decimal notation
	// let's forge a failure message

	// Let's figure out the notation
	switch {
	case strings.HasPrefix(currentValue, "0b"):
		return cleanedValue, errTimeDateFoundBinary
	case strings.HasPrefix(currentValue, "0x"):
		return cleanedValue, errTimeDateFoundHexadecimal
	case strings.HasPrefix(currentValue, "00") && len(currentValue) > 2:

		if parsedValue > 7 {
			strippedValue := strings.TrimLeft(currentValue, "0")

			return cleanedValue, fmt.Errorf("%w: %s (decimal) != %s (octal)", errTimeDateFoundOctalPaddingZero, strippedValue, cleanedValue)
		}

		// 00042 is a valid octal notation, but it equals to 34
		return cleanedValue, errTimeDateFoundOctalPaddingZero
	case strings.HasPrefix(currentValue, "0"):
		// this matches both "0" and "0o" octal notation
		return cleanedValue, errTimeDateFoundOctalPaddingZero
	}

	// anything else uses decimal notation
	// but let's validate it to be sure

	// convert back the number to a string, and compare it with the original one
	formattedValue := strconv.FormatInt(parsedValue, 10)
	if formattedValue != currentValue {
		// this can catch some edge cases like: 1_0 ...
		return cleanedValue, errTimeDateFoundAlternative
	}

	return bl.Value, nil
}
