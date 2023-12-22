package fixtures

import (
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"crypto/md5"
	md5 "crypto/md5"
	"strings"
	str "strings"     // MATCH /Import alias "str" is redundant/
	strings "strings" // MATCH /Import alias "strings" is redundant/
)

func UseRunpb() {
	runpb.RegisterTasksServer()
}

func UseMd5() {
	fmt.PrintLn(md5.Size)
}
