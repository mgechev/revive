package fixtures

import (
	b "bytes"
	"fmt"
	"net"
	"os"
)

func unhandledError1(a int) (int, error) {
	return a, nil
}

func prefixunhandledError1suffix(a int) (int, error) {
	return a, nil
}

func unhandledError2() error {
	// unhandledError1
	_, err := unhandledError1(1)
	unhandledError1(1)             // ignore
	prefixunhandledError1suffix(1) // MATCH /Unhandled error in call to function prefixunhandledError1suffix/
	return err
}

func testCase1() {
	// fmt\.Print
	fmt.Print(nil)        // ignore
	fmt.Println("")       // MATCH /Unhandled error in call to function fmt.Println/
	fmt.Printf("%d", 100) // MATCH /Unhandled error in call to function fmt.Printf/
	fmt.Fprintf(nil, "")  // MATCH /Unhandled error in call to function fmt.Fprintf/
}

func testCase2() {
	// os\.(Create|WriteFile|Chmod)
	os.Create("test")                                             // ignore
	os.Chmod("test_file", os.ModeAppend)                          // ignore
	os.WriteFile("test_file", []byte("some data"), os.ModeAppend) // ignore

	os.WriteFile("test_file", []byte("some data"), os.ModeAppend) // ignore

	_ = os.Chdir("..")
	os.Chdir("..") // MATCH /Unhandled error in call to function os.Chdir/
}

func testCase3() {
	// net\..*
	net.Dial("tcp", "127.0.0.1")                 // ignore
	net.ResolveTCPAddr("tcp4", "localhost:8080") // ignore
}

func testCase4() {
	// bytes\.Buffer\.Write
	b1 := b.Buffer{}
	b2 := &b.Buffer{}
	b1.Write(nil) // ignore
	b2.Write(nil) // ignore

	b2.Read([]byte("bytes")) // MATCH /Unhandled error in call to function bytes.Buffer.Read/
}

type unhandledErrorStruct1 struct {
}

func (s unhandledErrorStruct1) reterr() error {
	return nil
}

type unhandledErrorStruct2 struct {
}

func (s unhandledErrorStruct2) reterr() error {
	return nil
}

func (s *unhandledErrorStruct2) reterr1() error {
	return nil
}

func testCase5() {
	// fixtures\.unhandledErrorStruct2\.reterr
	s1 := unhandledErrorStruct1{}
	_ = s1.reterr()
	s1.reterr() // MATCH /Unhandled error in call to function fixtures.unhandledErrorStruct1.reterr/

	s2 := unhandledErrorStruct2{}
	s2.reterr()  // ignore
	s2.reterr1() // MATCH /Unhandled error in call to function fixtures.unhandledErrorStruct2.reterr1/
}
