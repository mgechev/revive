// Test that dot imports are flagged.

package fixtures

import (
	. "context" // in allowedPackages (standard library => just the name without full path)
	. "errors"  // in allowedPackages (standard library => just the name without full path)
	. "fmt"     // MATCH /should not use dot imports/
	"math/rand"
	tmplt "text/template"

	. "github.com/BurntSushi/toml" // in allowedPackages (not in the standard library)
)

var _ Stringer // from "fmt"
var _ = New("fake error")
var _ = Background()

var _ = Position{} // check a package not in the standard library

var _ = rand.Rand{}      // check non-alias package
var _ = tmplt.Template{} // check alias package
