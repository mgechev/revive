// Test of confusing-naming rule.

// Package pkg ...
package pkg

type foo struct {}

func (t foo) aFoo() { 
	return
} 



func (t *foo) AFoo() { // MATCH /Method 'AFoo' differs only by capitalization to method 'aFoo'/
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

func AGlobal(){ // MATCH /Method 'AGlobal' differs only by capitalization to method 'aGlobal'/
}

func ABar() { // Should not warn 

}

func aFoo() { // Should not warn 

}

func (t foo) ABar() { // Should not warn 
	return
} 

func (t bar) ABar() { // MATCH /Method 'ABar' differs only by capitalization to method 'aBar'/
	return
}