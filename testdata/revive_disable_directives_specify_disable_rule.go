package fixtures

// Disable rule is specified

//revive:disable:exported
func ExportedDisableRule1() {
}

//revive:enable

//revive:disable-next-line:exported
func ExportedDisableRule2() {
}

//revive:enable

// Disable rule is not specified

// MATCH:21 /rule name for lint disabling not found/

//revive:disable
func ExportedDisableRule3() { // MATCH /exported function ExportedDisableRule3 should have comment or be unexported/
}

//revive:enable

// MATCH:29 /rule name for lint disabling not found/

//revive:disable-next-line
func ExportedDisableRule4() { // MATCH /exported function ExportedDisableRule4 should have comment or be unexported/
}

//revive:enable
