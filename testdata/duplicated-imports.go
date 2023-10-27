package fixtures

import(
	"crypto/md5"
        "strings"
        _ "crypto/md5" // MATCH /Package "crypto/md5" already imported/
        str "strings"  // MATCH /Package "strings" already imported/
)
