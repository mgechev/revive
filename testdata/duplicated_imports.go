package fixtures

import (
	"crypto/md5"
	_ "crypto/md5" // MATCH /Package "crypto/md5" already imported/
	"strings"
	str "strings" // MATCH /Package "strings" already imported/
)
