// Test that blank imports in library packages are flagged.

// Package foo ...
package foo

// The instructions need to go before the imports below so they will not be
// mistaken for documentation.

import _ "encoding/json"

/* MATCH:9 /a blank import should be only in a main or test package, or have a comment justifying it/ */

import (
	"fmt"

	_ "os"
	/* MATCH:16 /a blank import should be only in a main or test package, or have a comment justifying it/ */

	_ "net/http"
	/* MATCH:19 /a blank import should be only in a main or test package, or have a comment justifying it/ */
	_ "path"
)

import _ "encoding/base64" // Don't gripe about this

import (
	// Don't gripe about these next two lines.
	_ "compress/zlib"

	_ "syscall"
	/* MATCH:30 /a blank import should be only in a main or test package, or have a comment justifying it/ */
	_ "path/filepath"
)

import (
	"go/ast"
	_ "go/scanner" // Don't gripe about this or the following line.
	_ "go/token"
)

import (
	_ "embed"
	/* MATCH:42 /a blank import should be only in a main or test package, or have a comment justifying it/ */
)

var (
	_ fmt.Stringer // for "fmt"
	_ ast.Node     // for "go/ast"
)
