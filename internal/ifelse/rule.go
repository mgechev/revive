package ifelse

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

type (
	Rule interface {
		CheckIfElse(chain Chain) (failMsg string)
	}
	Target int
)

// Apply evaluates the given Rule on if-else chains found within the given AST,
// and returns the failures.
//
// Note that in if-else chain with multiple "if" blocks, only the *last* one is checked,
// that is to say, given:
//
//	if foo {
//	    ...
//	} else if bar {
//		...
//	} else {
//		...
//	}
//
// Only the block following "bar" is linted. This is because the rules that use this function
// do not presently have anything to say about earlier blocks in the chain.
func Apply(rule Rule, node ast.Node, target Target) []lint.Failure {
	v := &visitor{rule: rule, target: target}
	ast.Walk(v, node)
	return v.failures
}

const (
	TargetIf   Target = iota // linter line-number will target the "if" statemenet
	TargetElse               // linter lint-number will target the "else" statmenet
)

type visitor struct {
	failures []lint.Failure
	target   Target
	rule     Rule
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	ifStmt, ok := node.(*ast.IfStmt)
	if !ok {
		return v
	}

	v.visitChain(ifStmt, Chain{})
	return nil
}

func (v *visitor) visitChain(ifStmt *ast.IfStmt, chain Chain) {
	// look for other if-else chains nested inside this if { } block
	ast.Walk(v, ifStmt.Body)

	if ifStmt.Else == nil {
		// no else branch
		return
	}

	if as, ok := ifStmt.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
		chain.HasIfInitializer = true
	}
	chain.IfTerminator = BlockTerminator(ifStmt.Body)

	switch elseBlock := ifStmt.Else.(type) {
	case *ast.IfStmt:
		if !chain.IfTerminator.DeviatesControlFlow() {
			chain.HasPriorNonReturn = true
		}
		v.visitChain(elseBlock, chain)
	case *ast.BlockStmt:
		// look for other if-else chains nested inside this else { } block
		ast.Walk(v, elseBlock)
		chain.ElseTerminator = BlockTerminator(elseBlock)
		if failMsg := v.rule.CheckIfElse(chain); failMsg != "" {
			if chain.HasIfInitializer {
				// if statement has a := initializer, so we might need to move the assignment
				// onto its own line in case the body references it
				failMsg += " (move short variable declaration to its own line if necessary)"
			}
			v.failures = append(v.failures, lint.Failure{
				Confidence: 1,
				Node:       v.targetNode(ifStmt),
				Failure:    failMsg,
			})
		}
	default:
		panic("invalid node type for else")
	}
}

func (v *visitor) targetNode(ifStmt *ast.IfStmt) ast.Node {
	switch v.target {
	case TargetIf:
		return ifStmt
	case TargetElse:
		return ifStmt.Else
	default:
		panic("bad target")
	}
}
