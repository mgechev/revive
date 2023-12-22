package fixtures

import (
	runpb "cloud.google.com/go/run/apiv2/runpb"
	"crypto/md5"
	md5 "crypto/md5"
	"strings"
	str "strings" 
	strings "strings" // MATCH /Import alias "strings" is redundant/
)

func UseRunpb() {
	runpb.RegisterTasksServer()
}

func UseMd5() {
	fmt.PrintLn(md5.Size)
}
