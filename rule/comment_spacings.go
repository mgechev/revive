// Package rule implements revive's linting rules.
package rule

import (
	"fmt"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

// CommentSpacingsRule check the whether there is a space between
// the comment symbol( // ) and the start of the comment text
type CommentSpacingsRule struct {
	allowList []string
	sync.Mutex
}

func (r *CommentSpacingsRule) configure(arguments lint.Arguments) error {
	r.Lock()
	defer r.Unlock()
	if r.allowList != nil {
		return nil // already configured
	}

	r.allowList = []string{}
	for _, arg := range arguments {
		allow, ok := arg.(string) // Alt. non panicking version
		if !ok {
			return fmt.Errorf("invalid argument %v for %s; expected string but got %T", arg, r.Name(), arg)
		}
		r.allowList = append(r.allowList, `//`+allow)
	}
	return nil
}

// Apply the rule.
func (r *CommentSpacingsRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	var failures []lint.Failure
	err := r.configure(arguments)
	if err != nil {
		return failures, err
	}

	for _, cg := range file.AST.Comments {
		for _, comment := range cg.List {
			commentLine := comment.Text
			if len(commentLine) < 3 {
				continue // nothing to do
			}

			isMultiLineComment := commentLine[1] == '*'
			isOK := commentLine[2] == '\n'
			if isMultiLineComment && isOK {
				continue
			}

			isOK = (commentLine[2] == ' ') || (commentLine[2] == '\t')
			if isOK {
				continue
			}

			if r.isAllowed(commentLine) {
				continue
			}

			failures = append(failures, lint.Failure{
				Node:       comment,
				Confidence: 1,
				Category:   "style",
				Failure:    "no space between comment delimiter and comment text",
			})
		}
	}
	return failures, nil
}

// Name yields this rule name.
func (*CommentSpacingsRule) Name() string {
	return "comment-spacings"
}

func (r *CommentSpacingsRule) isAllowed(line string) bool {
	for _, allow := range r.allowList {
		if strings.HasPrefix(line, allow) {
			return true
		}
	}

	return isDirectiveComment(line)
}
