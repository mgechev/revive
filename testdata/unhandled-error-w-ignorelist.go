package fixtures

import (
	"fmt"
	"net"
	"os"
)

func unhandledError1(a int) (int, error) {
	return a, nil
}

func unhandledError2() error {
	_, err := unhandledError1(1)
	unhandledError1(1)
	prefixunhandledError1suffix(1)
	fmt.Fprintf(nil, "") // MATCH /Unhandled error in call to function fmt.Fprintf/
	net.Dial("tcp", "127.0.0.1")
	net.ResolveTCPAddr("tcp4", "localhost:8080")
	os.Chdir("..")
	_ = os.Chdir("..")
	return err
}
