package fixtures

import (
	"crypto/md5"
	md5 "crypto/md5" // MATCH /Import alias "md5" is redundant/

	runpb "cloud.google.com/go/run/apiv2/runpb" // MATCH /Import alias "runpb" is redundant/
)

func UseRunpb() {
	runpb.RegisterTasksServer()
}

func UseMd5() {
	fmt.PrintLn(md5.Size)
}
