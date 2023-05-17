package ifelse

import (
	"fmt"
	"go/ast"
	"go/token"
)

// Branch contains information about a branch within an if-else chain.
type Branch struct {
	BranchKind
	Call // The function called at the end for kind Panic or Exit.
}

// BlockBranch gets the Branch of an ast.BlockStmt.
func BlockBranch(block *ast.BlockStmt) Branch {
	blockLen := len(block.List)
	if blockLen == 0 {
		return Branch{BranchKind: Empty}
	}

	switch stmt := block.List[blockLen-1].(type) {
	case *ast.ReturnStmt:
		return Branch{BranchKind: Return}
	case *ast.BlockStmt:
		return BlockBranch(stmt)
	case *ast.BranchStmt:
		switch stmt.Tok {
		case token.BREAK:
			return Branch{BranchKind: Break}
		case token.CONTINUE:
			return Branch{BranchKind: Continue}
		case token.GOTO:
			return Branch{BranchKind: Goto}
		}
	case *ast.ExprStmt:
		fn, ok := ExprCall(stmt)
		if !ok {
			break
		}
		kind, ok := DeviatingFuncs[fn]
		if ok {
			return Branch{BranchKind: kind, Call: fn}
		}
	}

	return Branch{BranchKind: Regular}
}

// String returns a brief string representation
func (b Branch) String() string {
	switch b.BranchKind {
	case Panic, Exit:
		return fmt.Sprintf("... %v()", b.Call)
	default:
		return b.BranchKind.String()
	}
}

// LongString returns a longer form string representation
func (b Branch) LongString() string {
	switch b.BranchKind {
	case Panic, Exit:
		return fmt.Sprintf("call to %v function", b.Call)
	default:
		return b.BranchKind.LongString()
	}
}
