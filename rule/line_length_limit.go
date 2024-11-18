package rule

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/token"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/mgechev/revive/lint"
)

// LineLengthLimitRule lints given else constructs.
type LineLengthLimitRule struct {
	max int

	configureOnce sync.Once
}

const defaultLineLengthLimit = 80

func (r *LineLengthLimitRule) configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.max = defaultLineLengthLimit
		return nil
	}

	maxLength, ok := arguments[0].(int64) // Alt. non panicking version
	if !ok || maxLength < 0 {
		return errors.New(`invalid value passed as argument number to the "line-length-limit" rule`)
	}

	r.max = int(maxLength)
	return nil
}

// Apply applies the rule to given file.
func (r *LineLengthLimitRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	var configErr error
	r.configureOnce.Do(func() { configErr = r.configure(arguments) })
	if configErr != nil {
		return nil, configErr
	}

	var failures []lint.Failure

	checker := lintLineLengthNum{
		max:  r.max,
		file: file,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	checker.check()

	return failures, nil
}

// Name returns the rule name.
func (*LineLengthLimitRule) Name() string {
	return "line-length-limit"
}

type lintLineLengthNum struct {
	max       int
	file      *lint.File
	onFailure func(lint.Failure)
}

func (r lintLineLengthNum) check() {
	f := bytes.NewReader(r.file.Content())
	spaces := strings.Repeat(" ", 4) // tab width = 4
	l := 1
	s := bufio.NewScanner(f)
	for s.Scan() {
		t := s.Text()
		t = strings.ReplaceAll(t, "\t", spaces)
		c := utf8.RuneCountInString(t)
		if c > r.max {
			r.onFailure(lint.Failure{
				Category: "code-style",
				Position: lint.FailurePosition{
					// Offset not set; it is non-trivial, and doesn't appear to be needed.
					Start: token.Position{
						Filename: r.file.Name,
						Line:     l,
						Column:   0,
					},
					End: token.Position{
						Filename: r.file.Name,
						Line:     l,
						Column:   c,
					},
				},
				Confidence: 1,
				Failure:    fmt.Sprintf("line is %d characters, out of limit %d", c, r.max),
			})
		}
		l++
	}
}
