package ifelse

// Kind is a classifier for if-else branches. It says whether the branch is empty,
// and whether the branch ends with a statement that can deviate control flow.
type Kind int

const (
	Empty Kind = iota
	Regular
	Return
	Continue
	Break
	Goto
	Panic
	Exit
)

func (k Kind) IsEmpty() bool  { return k == Empty }
func (k Kind) IsReturn() bool { return k == Return }

// DeviatesControlFlow returns true if the program follows regular control flow at
// branch termination (that is to say, control flows to the first statement following
// the if-else chain).
func (k Kind) DeviatesControlFlow() bool {
	switch k {
	case Empty, Regular:
		return false
	case Return, Continue, Break, Goto, Panic, Exit:
		return true
	default:
		panic("invalid kind")
	}
}

func (k Kind) String() string {
	switch k {
	case Empty:
		return ""
	case Regular:
		return "..."
	case Return:
		return "... return"
	case Continue:
		return "... continue"
	case Break:
		return "... break"
	case Goto:
		return "... goto"
	case Panic:
		return "... panic()"
	case Exit:
		return "... os.Exit()"
	default:
		panic("invalid kind")
	}
}

func (k Kind) LongString() string {
	switch k {
	case Empty:
		return "an empty block"
	case Regular:
		return "a regular statement"
	case Return:
		return "a return statement"
	case Continue:
		return "a continue statement"
	case Break:
		return "a break statement"
	case Goto:
		return "a goto statement"
	case Panic:
		return "a function call that panics"
	case Exit:
		return "a function call that exits the program"
	default:
		panic("invalid kind")
	}
}
