// Test for name linting.

// Package pkg_with_underscores ...
package pkg_with_underscores // MATCH /don't use an underscore in package name/

import (
	"io"
	"net"
	net_http "net/http" // renamed deliberately
	"net/url"
)

import "C"

var safeUrl = "HttPS://..iaMHost..Test:443/paTh^A%ef//./%41PaTH/..//?" // MATCH /var safeUrl should be safeURL/
var var_name int                                                       // MATCH /don't use underscores in Go names; var var_name should be varName/

type t_wow struct { // MATCH /don't use underscores in Go names; type t_wow should be tWow/
	x_damn int      // MATCH /don't use underscores in Go names; struct field x_damn should be xDamn/
	Url    *url.URL // MATCH /struct field Url should be URL/
}

const fooId = "blah" // MATCH /const fooId should be fooID/

func f_it() { // MATCH /don't use underscores in Go names; func f_it should be fIt/
	more_underscore := 4 // MATCH /don't use underscores in Go names; var more_underscore should be moreUnderscore/
	_ = more_underscore
	var err error
	if isEof := (err == io.EOF); isEof { // MATCH /var isEof should be isEOF/
		more_underscore = 7 // should be okay
	}

	x := net_http.Request{} // should be okay
	_ = x

	var ips []net.IP
	for _, theIp := range ips { // MATCH /range var theIp should be theIP/
		_ = theIp
	}

	switch myJson := g(); { // MATCH /var myJson should be myJSON/
	default:
		_ = myJson
	}
	var y net_http.ResponseWriter // an interface
	switch tApi := y.(type) {     // MATCH /var tApi should be tAPI/
	default:
		_ = tApi
	}

	var c chan int
	select {
	case qId := <-c: // MATCH /var qId should be qID/
		_ = qId
	}
}

// Common styles in other languages that don't belong in Go.
const (
	CPP_CONST = 1 // MATCH /don't use ALL_CAPS in Go names; use CamelCase/

	HTML  = 3 // okay; no underscore
	X509B = 4 // ditto
)

func f(bad_name int)                    {}            // MATCH /don't use underscores in Go names; func parameter bad_name should be badName/
func g() (no_way int)                   { return 0 }  // MATCH /don't use underscores in Go names; func result no_way should be noWay/
func (t *t_wow) f(more_under string)    {}            // MATCH /don't use underscores in Go names; method parameter more_under should be moreUnder/
func (t *t_wow) g() (still_more string) { return "" } // MATCH /don't use underscores in Go names; method result still_more should be stillMore/

type i interface {
	CheckHtml() string // okay; interface method names are often constrained by the concrete types' method names

	F(foo_bar int) // MATCH /don't use underscores in Go names; interface method parameter foo_bar should be fooBar/
}

// All okay; underscore between digits
const case1_1 = 1

type case2_1 struct {
	case2_2 int
}

func case3_1(case3_2 int) (case3_3 string) {
	case3_4 := 4
	_ = case3_4

	return ""
}

type t struct{}

func (t) LastInsertId() (int64, error) { return 0, nil } // okay because it matches a known style violation

//export exported_to_c
func exported_to_c() {} // okay: https://github.com/golang/lint/issues/144

//export exported_to_c_with_arg
func exported_to_c_with_arg(but_use_go_param_names int) // MATCH /don't use underscores in Go names; func parameter but_use_go_param_names should be butUseGoParamNames/

// This is an exported C function with a leading doc comment.
//
//export exported_to_c_with_comment
func exported_to_c_with_comment() {} // okay: https://github.com/golang/lint/issues/144

//export maybe_exported_to_CPlusPlusWithCamelCase
func maybe_exported_to_CPlusPlusWithCamelCase() {} // okay: https://github.com/golang/lint/issues/144

// WhyAreYouUsingCapitalLetters_InACFunctionName is a Go-exported function that
// is also exported to C as a name with underscores.
//
// Don't do that. If you want to use a C-style name for a C export, make it
// lower-case and leave it out of the Go-exported API.
//
//export WhyAreYouUsingCapitalLetters_InACFunctionName
func WhyAreYouUsingCapitalLetters_InACFunctionName() {} // MATCH /don't use underscores in Go names; func WhyAreYouUsingCapitalLetters_InACFunctionName should be WhyAreYouUsingCapitalLettersInACFunctionName/
