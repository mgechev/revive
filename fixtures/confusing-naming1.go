// Test of confusing-naming rule.

// Package pkg ...
package pkg

type foo struct {}

func (t foo) aFoo() { 
	return
} 



func (t *foo) AFoo() { // MATCH /Method 'AFoo' differs only by capitalization to other method of 'pkg/foo'/
	return
}
 
type bar struct {}

func (t *bar) aBar() { 
	return
} 

func (t *bar) aFoo() { // Should not warn 
	return
}



func aGlobal(){

}

func AGlobal(){ // MATCH /Function 'AGlobal' differs only by capitalization to other function/
}

func ABar() { // Should not warn 

}

func aFoo() { // Should not warn 

}

func (t foo) ABar() { // Should not warn 
	return
} 

func (t bar) ABar() { // MATCH /Method 'ABar' differs only by capitalization to other method of 'pkg/bar'/
	return
}