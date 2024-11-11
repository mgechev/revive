package fixtures

func getfoo() (int, int, error) { // MATCH /unnamed results of the same type may be confusing, consider using named results/

}

func getBar(a, b int) (int, error, int) {
}

func Getbaz(a string, b int) (int, float32, string, string) { // MATCH /unnamed results of the same type may be confusing, consider using named results/

}

func GetTaz(a string, b int) string {

}

func (t *t) GetTaz(a int, b int)  {

}
