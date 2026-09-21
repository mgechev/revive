package fixtures

//revive:disable
func Public1() {
}

//revive:enable

func Public2() { // MATCH /exported function Public2 should have comment or be unexported/
}

//revive:disable:exported
func Public3() {
}

//revive:enable:exported

//revive:disable:random

func Public4() { // MATCH /exported function Public4 should have comment or be unexported/
}

//revive:enable:random
