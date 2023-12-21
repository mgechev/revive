package fixtures

import (
	runpb "cloud.google.com/go/run/apiv2/runpb" // MATCH /Import alias "runpb" is redundant/

	md5 "crypto/md5" // MATCH /Import alias "md5" is redundant/

	strings "strings" // MATCH /Import alias "strings" is redundant/

	"crypto/md5"
	_ "crypto/md5" // MATCH /Import alias "_" is redundant/

	"strings"
	str "strings" // MATCH /Import alias "str" is redundant/

)

func UseRunpb() {
	runpb.RegisterTasksServer()
}

func UseMd5() {
	fmt.PrintLn(md5.Size)
}
