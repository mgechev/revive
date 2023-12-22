package fixtures

import(
	"crypto/md5"
		"strings"
		_ "crypto/md5"
		str "strings"
		strings "strings"  // MATCH /Import alias "strings" is redundant/
		crypto "crypto/md5"
        redundant "abc/redundant" // MATCH /Import alias "redundant" is redundant/
		md5 "crypto/md5"  // MATCH /Import alias "md5" is redundant/
)

