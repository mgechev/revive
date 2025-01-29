// Package fix_tures is a testing package
// tests for issue #540
//revive:disable-next-line:var-naming
package fix_tures

func foo1() {
	// something before the annotation
	//revive:disable-next-line:var-naming
	// something after
	var invalid_name = 0
}

//revive:disable-next-line:var-naming
func (source Source) BaseApiURL() string {}
