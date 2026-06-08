package fixtures

// Disable reason and rule are specified

//revive:disable:exported reason for disabling rule
func ExportedDisableReasonRule1() {
}

//revive:enable

//revive:disable-next-line:exported reason for disabling rule
func ExportedDisableReasonRule2() {
}

//revive:enable

// Disable reason and rule are not specified

// MATCH:21 /reason of lint disabling not found/

//revive:disable
func ExportedDisableReasonRule3() { // MATCH /exported function ExportedDisableReasonRule3 should have comment or be unexported/
}

//revive:enable

// MATCH:29 /reason of lint disabling not found/

//revive:disable-next-line
func ExportedDisableReasonRule4() { // MATCH /exported function ExportedDisableReasonRule4 should have comment or be unexported/
}

//revive:enable
