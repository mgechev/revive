// Package rule implements revive's linting rules.
package rule

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"sync"

	"github.com/mgechev/revive/lint"
)

// FileLengthLimitRule lints the number of lines in a file.
type FileLengthLimitRule struct {
	// max is the maximum number of lines allowed in a file. 0 means the rule is disabled.
	max int
	// skipComments indicates whether to skip comment lines when counting lines.
	skipComments bool
	// skipBlankLines indicates whether to skip blank lines when counting lines.
	skipBlankLines bool
	sync.Mutex
}

// Apply applies the rule to given file.
func (r *FileLengthLimitRule) Apply(file *lint.File, arguments lint.Arguments) ([]lint.Failure, error) {
	var failures []lint.Failure
	err := r.configure(arguments)
	if err != nil {
		return failures, err
	}

	if r.max <= 0 {
		// when max is negative or 0 the rule is disabled
		return nil, nil
	}

	all := 0
	blank := 0
	scanner := bufio.NewScanner(bytes.NewReader(file.Content()))
	for scanner.Scan() {
		all++
		if len(bytes.TrimSpace(scanner.Bytes())) == 0 {
			blank++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	lines := all
	if r.skipComments {
		lines -= countCommentLines(file.AST.Comments)
	}

	if r.skipBlankLines {
		lines -= blank
	}

	if lines <= r.max {
		return nil, nil
	}

	return []lint.Failure{
		{
			Category:   "code-style",
			Confidence: 1,
			Position: lint.FailurePosition{
				Start: token.Position{
					Filename: file.Name,
					Line:     all,
				},
			},
			Failure: fmt.Sprintf("file length is %d lines, which exceeds the limit of %d", lines, r.max),
		},
	}, nil
}

func (r *FileLengthLimitRule) configure(arguments lint.Arguments) error {
	r.Lock()
	defer r.Unlock()

	if r.max != 0 {
		return nil // already configured
	}

	if len(arguments) < 1 {
		return nil // use default
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "file-length-limit" rule. Expecting a k,v map, got %T`, arguments[0])
	}
	for k, v := range argKV {
		switch k {
		case "max":
			maxLines, ok := v.(int64)
			if !ok || maxLines < 0 {
				return fmt.Errorf(`invalid configuration value for max lines in "file-length-limit" rule; need positive int64 but got %T`, arguments[0])
			}
			r.max = int(maxLines)
		case "skipComments":
			skipComments, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skip comments in "file-length-limit" rule; need bool but got %T`, arguments[1])
			}
			r.skipComments = skipComments
		case "skipBlankLines":
			skipBlankLines, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for skip blank lines in "file-length-limit" rule; need bool but got %T`, arguments[2])
			}
			r.skipBlankLines = skipBlankLines
		}
	}
	return nil
}

// Name returns the rule name.
func (*FileLengthLimitRule) Name() string {
	return "file-length-limit"
}

func countCommentLines(comments []*ast.CommentGroup) int {
	count := 0
	for _, cg := range comments {
		for _, comment := range cg.List {
			if len(comment.Text) < 2 {
				continue
			}
			switch comment.Text[1] {
			case '/': // single-line comment
				count++
			case '*': // multi-line comment
				count += strings.Count(comment.Text, "\n") + 1
			}
		}
	}
	return count
}
