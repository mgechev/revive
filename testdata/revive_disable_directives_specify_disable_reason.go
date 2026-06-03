package fixtures

// Disable reason is specified

//revive:disable // it's ok to have exported function without comment
func Exported1() {
}

//revive:enable

//revive:disable:exported // it's ok to have exported function without comment
func Exported2() {
}

//revive:enable

//revive:disable-next-line // it's ok to have exported function without comment
func Exported3() {
}

//revive:disable-next-line:exported // it's ok to have exported function without comment
func Exported4() {
}

//revive:enable

// Disable reason is not specified

// MATCH:31 /reason of lint disabling not found/

//revive:disable
func Exported5() { // MATCH /exported function Exported5 should have comment or be unexported/
}

//revive:enable

// MATCH:39 /reason of lint disabling not found/

//revive:disable:exported
func Exported6() { // MATCH /exported function Exported6 should have comment or be unexported/
}

//revive:enable

// MATCH:47 /reason of lint disabling not found/

//revive:disable-next-line
func Exported7() { // MATCH /exported function Exported7 should have comment or be unexported/
}

//revive:enable

// MATCH:55 /reason of lint disabling not found/

//revive:disable-next-line:exported
func Exported8() { // MATCH /exported function Exported8 should have comment or be unexported/
}

//revive:enable
