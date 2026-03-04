// From https://github.com/fatih/structtag/blob/v1.2.0/tags.go
/*
Copyright (c) 2017, Fatih Arslan
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of structtag nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

// Package structtag provides utilities to parse [struct] tags.
//
// [struct]: https://go.dev/ref/spec#Struct_types
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
		// Skip leading ASCII whitespace to match scanKey's delimiter rules.
		for tag != "" && tag[0] <= ' ' {
			tag = tag[1:]
		}
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
// The tag string must start with a double-quote character.
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
