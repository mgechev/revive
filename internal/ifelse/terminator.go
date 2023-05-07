package ifelse

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Terminator contains information about the end of a branch within an if-else chain.
type Terminator struct {
	Kind
	Call // The function called at the end for kind Panic or Exit.
}

// BlockTerminator gets the Terminator of an ast.BlockStmt.
func BlockTerminator(block *ast.BlockStmt) Terminator {
	blockLen := len(block.List)
	if blockLen == 0 {
		return Terminator{Kind: Empty}
	}

	switch stmt := block.List[blockLen-1].(type) {
	case *ast.ReturnStmt:
		return Terminator{Kind: Return}
	case *ast.BlockStmt:
		return BlockTerminator(stmt)
	case *ast.BranchStmt:
		switch stmt.Tok {
		case token.BREAK:
			return Terminator{Kind: Break}
		case token.CONTINUE:
			return Terminator{Kind: Continue}
		case token.GOTO:
			return Terminator{Kind: Goto}
		}
	case *ast.ExprStmt:
		fn, ok := ExprCall(stmt)
		if !ok {
			break
		}
		kind, ok := DeviatingCalls[fn]
		if ok {
			return Terminator{Kind: kind, Call: fn}
		}
	}

	return Terminator{Kind: Regular}
}

func (b Terminator) String() string {
	switch b.Kind {
	case Panic, Exit:
		return fmt.Sprintf("... %v()", b.Call)
	default:
		return b.Kind.String()
	}
}

func (b Terminator) LongString() string {
	switch b.Kind {
	case Panic, Exit:
		return fmt.Sprintf("call to %v function", b.Call)
	default:
		return b.Kind.LongString()
	}
}
