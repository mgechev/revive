package fixtures

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
