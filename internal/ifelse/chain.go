package ifelse

// Chain contains information about an if-else chain.
type Chain struct {
	IfTerminator      Terminator // what happens at the end of the "if" block
	ElseTerminator    Terminator // what happens at the end of the "else" block
	HasIfInitializer  bool       // is there an "if"-initializer somewhere in the chain?
	HasPriorNonReturn bool       // is there a prior "if" block that does NOT deviate control flow?
}
