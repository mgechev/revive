package rule

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/logging"
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
	// timeDateArgumentNames are the names of the arguments of time.Date
	timeDateArgumentNames = []string{
		"year",
		"month",
		"day",
		"hour",
		"minute",
		"second",
		"nanosecond",
		"timezone",
	}

	// timeDateArity is the number of arguments of time.Date
	timeDateArity = len(timeDateArgumentNames)
)

var timeDateArgumentBoundaries = map[string][2]int64{
	"month":      {1, 12},
	"day":        {1, 31},
	"hour":       {0, 23},
	"minute":     {0, 59},
	"second":     {0, 59},
	"nanosecond": {0, 999999999},
}

func (w lintTimeDate) Visit(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if !ok || len(ce.Args) != timeDateArity {
		return w
	}
	if !isPkgDot(ce.Fun, "time", "Date") {
		return w
	}

	// The last argument is a timezone, there is no need to check it, also it has a different type
	for pos, realArg := range ce.Args[:timeDateArity-1] {
		arg := realArg
		var isNegative bool
		ua, found := arg.(*ast.UnaryExpr)
		if found {
			// this is the argument we want to check
			arg = ua.X
			if ua.Op == token.SUB {
				isNegative = true
			}
		}

		bl, ok := arg.(*ast.BasicLit)
		if !ok {
			continue
		}
		fieldName := timeDateArgumentNames[pos]

		parsedValue, err := parseDecimalInteger(bl)
		if err == nil {
			if isNegative {
				parsedValue = -parsedValue
			}

			boundaries, found := timeDateArgumentBoundaries[fieldName]
			if found {
				if parsedValue < boundaries[0] || parsedValue > boundaries[1] {
					w.onFailure(lint.Failure{
						Category:   "time",
						Node:       realArg,
						Confidence: 0.8,
						Failure: fmt.Sprintf(
							"time.Date %s argument is supposed to be between %d and %d, found: %d",
							fieldName, boundaries[0], boundaries[1], parsedValue),
					})
				}
			}

			continue
		}

		if errors.Is(err, errParsedInvalid) {
			// This is not supposed to happen, let's be defensive
			// log the error, but continue

			logger, errLogger := logging.GetLogger()
			if errLogger != nil {
				// This is not supposed to happen, discard both errors
				continue
			}
			logger.With(
				"value", bl.Value,
				"kind", bl.Kind,
				"error", err.Error(),
			).Error("failed to parse time.Date argument")

			continue
		}

		confidence := 0.8 // default confidence
		errMessage := err.Error()
		replacedValue := strconv.FormatInt(parsedValue, 10)
		instructions := fmt.Sprintf("use %s instead of %s", replacedValue, bl.Value)
		switch {
		case errors.Is(err, errParsedOctalWithZero):
			// people can use 00, 01, 02, 03, 04, 05, 06, and 07 if they want.
			confidence = 0.5

		case errors.Is(err, errParsedOctalWithPaddingZeroes):
			// This is a clear mistake.
			// example with 000123456 (octal) is about 123456 or 42798 ?
			confidence = 1

			strippedValue := strings.TrimLeft(bl.Value, "0")
			if strippedValue == "" {
				// avoid issue with 00000000
				strippedValue = "0"
			}

			if strippedValue != replacedValue {
				instructions = fmt.Sprintf(
					"choose between %s and %s (decimal value of %s octal value)",
					strippedValue, replacedValue, strippedValue,
				)
			}
		}

		w.onFailure(lint.Failure{
			Category:   "time",
			Node:       bl,
			Confidence: confidence,
			Failure: fmt.Sprintf(
				"use decimal digits for time.Date %s argument: %s found: %s",
				timeDateArgumentNames[pos], errMessage, instructions),
		})
	}

	return w
}

var (
	errParsedOctal                  = errors.New("octal notation")
	errParsedOctalWithZero          = errors.New("octal notation with leading zero")
	errParsedOctalWithPaddingZeroes = errors.New("octal notation with padding zeroes")
	errParsedHexadecimal            = errors.New("hexadecimal notation")
	errParseBinary                  = errors.New("binary notation")
	errParsedFloat                  = errors.New("float literal")
	errParsedExponential            = errors.New("exponential notation")
	errParsedAlternative            = errors.New("alternative notation")
	errParsedInvalid                = errors.New("invalid notation")
)

func parseDecimalInteger(bl *ast.BasicLit) (int64, error) {
	currentValue := strings.ToLower(bl.Value)

	if currentValue == "0" {
		// skip 0 as it is a valid value for all the arguments
		return 0, nil
	}

	switch bl.Kind {
	case token.FLOAT:
		// someone used a float literal, while they should have used an integer literal.
		parsedValue, err := strconv.ParseFloat(currentValue, 64)
		if err != nil {
			// This is not supposed to happen
			return 0, fmt.Errorf(
				"%w: %s: %w",
				errParsedInvalid,
				"failed to parse number as float",
				err,
			)
		}

		// this will convert back the number to a string
		if strings.Contains(currentValue, "e") {
			return int64(parsedValue), errParsedExponential
		}

		return int64(parsedValue), errParsedFloat

	case token.INT:
		// we expect this format

	default:
		// This is not supposed to happen
		return 0, fmt.Errorf(
			"%w: %s",
			errParsedInvalid,
			"unexpected kind of literal",
		)
	}

	// Parse the number with base=0 that allows to accept all number formats and base
	parsedValue, err := strconv.ParseInt(currentValue, 0, 64)
	if err != nil {
		// This is not supposed to happen
		return 0, fmt.Errorf(
			"%w: %s: %w",
			errParsedInvalid,
			"failed to parse number as integer",
			err,
		)
	}

	// Let's figure out the notation to return an error
	switch {
	case strings.HasPrefix(currentValue, "0b"):
		return parsedValue, errParseBinary
	case strings.HasPrefix(currentValue, "0x"):
		return parsedValue, errParsedHexadecimal
	case strings.HasPrefix(currentValue, "0"):
		// this matches both "0" and "0o" octal notation.

		switch currentValue {
		// people can use 00, 01, 02, 03, 04, 05, 06, 07, if they want
		case "00", "01", "02", "03", "04", "05", "06", "07":
			return parsedValue, errParsedOctalWithZero
		}

		if strings.HasPrefix(currentValue, "00") {
			// 00123456 (octal) is about 123456 or 42798 ?
			return parsedValue, errParsedOctalWithPaddingZeroes
		}

		return parsedValue, errParsedOctal
	}

	// Convert back the number to a string, and compare it with the original one
	formattedValue := strconv.FormatInt(parsedValue, 10)
	if formattedValue != currentValue {
		// This can catch some edge cases like: 1_0 ...
		return parsedValue, errParsedAlternative
	}

	return parsedValue, nil
}
