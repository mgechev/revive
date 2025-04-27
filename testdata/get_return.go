package fixtures

import "net/http"

func getfoo() {

}

func getBar(a, b int) { // MATCH /function 'getBar' seems to be a getter but it does not return any result/

}

func Getbaz(a string, b int) {

}

func GetTaz(a string, b int) string {

}

func (t *t) GetTaz(a string, b int) string {

}

func (t *t) GetSaz(a string, b int) { // MATCH /function 'GetSaz' seems to be a getter but it does not return any result/

}

func GetQux(a string, b int, c int, d string, e int64) { // MATCH /function 'GetQux' seems to be a getter but it does not return any result/

}

// non-regression test issue #1323
func (b *t) GetInfo(w http.ResponseWriter, r *http.Request) {}

func GetSomething(w http.ResponseWriter, r *http.Request, p int) {}

func GetSomethingElse(p int, w http.ResponseWriter, r *http.Request) {} // MATCH /function 'GetSomethingElse' seems to be a getter but it does not return any result/
