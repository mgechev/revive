// Package structtag provides utilities to parse and manipulate struct tags.
package structtag

import (
	"errors"
	"strconv"
	"strings"
)

var errValueSyntax = errors.New("invalid syntax for struct tag value")

// Tag defines a single struct's string literal tag.
type Tag struct {
	// Key is the tag key, such as json, xml, etc.
	// i.e., `json:"foo,omitempty"`. Here Key is "json".
	Key string

	// Name is the tag name.
	// i.e., `json:"foo,omitempty"`. Here Name is "foo".
	Name string

	// Options are the tag options.
	// i.e., `json:"foo,omitempty"`. Here Options is ["omitempty"].
	Options []string
}

// Parse parses a single struct field tag and returns the set of tags.
func Parse(tag string) ([]*Tag, error) {
	if tag == "" {
		return nil, nil
	}

	var tags []*Tag

	for tag != "" {
		// Skip leading space.
		tag = strings.TrimLeft(tag, " ")
		if tag == "" {
			break
		}

		i := scanKey(tag)
		if i == 0 {
			return nil, errors.New("invalid syntax for struct tag key")
		}
		if i+1 >= len(tag) || tag[i] != ':' {
			return nil, errors.New("invalid syntax for struct tag pair")
		}
		if tag[i+1] != '"' {
			return nil, errValueSyntax
		}

		key := tag[:i]
		tag = tag[i+1:]

		i, qvalue, err := scanValue(tag)
		if err != nil {
			return nil, err
		}

		value, err := strconv.Unquote(qvalue)
		if err != nil {
			return nil, errValueSyntax
		}

		name, options := parseTagValue(value)
		tags = append(tags, &Tag{
			Key:     key,
			Name:    name,
			Options: options,
		})

		tag = tag[i:]
	}

	if len(tags) == 0 {
		return nil, nil
	}

	return tags, nil
}

// scanKey finds the position of the colon that ends the tag key.
// It returns the index of the first invalid character (space, colon, quote, or control character),
// or len(tag) if all characters are valid.
func scanKey(tag string) int {
	for i := 0; i < len(tag); i++ {
		if tag[i] <= ' ' || tag[i] == ':' || tag[i] == '"' || tag[i] == 0x7f {
			return i
		}
	}
	return len(tag)
}

// scanValue scans a quoted string value and returns its index and quoted content.
func scanValue(tag string) (idx int, qvalue string, err error) {
	// Find closing quote, handling escapes.
	i := 1
	for i < len(tag) && tag[i] != '"' {
		if tag[i] == '\\' {
			i++
		}
		i++
	}
	if i >= len(tag) {
		return 0, "", errValueSyntax
	}
	return i + 1, tag[:i+1], nil
}

// parseTagValue parses an unquoted tag value into name and options.
// The format is "name" or "name,opt1,opt2,...".
func parseTagValue(value string) (name string, options []string) {
	parts := strings.Split(value, ",")
	name = parts[0]
	if len(parts) > 1 {
		options = parts[1:]
	}
	return name, options
}
