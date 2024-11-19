package ifelse

import (
	"go/ast"
	"go/token"

	"github.com/mgechev/revive/lint"
)

// Rule is an interface for linters operating on if-else chains
type Rule interface {
	CheckIfElse(chain Chain, args Args) (failMsg string)
}

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
func Apply(rule Rule, node ast.Node, target Target, args lint.Arguments) []lint.Failure {
	v := &visitor{rule: rule, target: target}
	for _, arg := range args {
		switch arg {
		case PreserveScope:
			v.args.PreserveScope = true
		case AllowJump:
			v.args.AllowJump = true
		}
	}
	ast.Walk(v, node)
	return v.failures
}

type visitor struct {
	failures []lint.Failure
	target   Target
	rule     Rule
	args     Args
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	switch stmt := node.(type) {
	case *ast.FuncDecl:
		if stmt.Body != nil {
			v.visitBlock(stmt.Body.List, Return)
		}
	case *ast.FuncLit:
		v.visitBlock(stmt.Body.List, Return)
	case *ast.ForStmt:
		v.visitBlock(stmt.Body.List, Continue)
	case *ast.RangeStmt:
		v.visitBlock(stmt.Body.List, Continue)
	case *ast.CaseClause:
		v.visitBlock(stmt.Body, Break)
	case *ast.BlockStmt:
		v.visitBlock(stmt.List, Regular)
	default:
		return v
	}
	return nil
}

func (v *visitor) visitBlock(stmts []ast.Stmt, endKind BranchKind) {
	for i, stmt := range stmts {
		ifStmt, ok := stmt.(*ast.IfStmt)
		if !ok {
			ast.Walk(v, stmt)
			continue
		}
		var chain Chain
		if i == len(stmts)-1 {
			chain.AtBlockEnd = true
			chain.BlockEndKind = endKind
		}
		v.visitIf(ifStmt, chain)
	}
}

func (v *visitor) visitIf(ifStmt *ast.IfStmt, chain Chain) {
	// look for other if-else chains nested inside this if { } block
	v.visitBlock(ifStmt.Body.List, chain.BlockEndKind)

	if as, ok := ifStmt.Init.(*ast.AssignStmt); ok && as.Tok == token.DEFINE {
		chain.HasInitializer = true
	}
	chain.If = BlockBranch(ifStmt.Body)

	if ifStmt.Else == nil {
		if v.args.AllowJump {
			v.checkRule(ifStmt, chain)
		}
		return
	}

	switch elseBlock := ifStmt.Else.(type) {
	case *ast.IfStmt:
		if !chain.If.Deviates() {
			chain.HasPriorNonDeviating = true
		}
		v.visitIf(elseBlock, chain)
	case *ast.BlockStmt:
		// look for other if-else chains nested inside this else { } block
		v.visitBlock(elseBlock.List, chain.BlockEndKind)

		chain.HasElse = true
		chain.Else = BlockBranch(elseBlock)
		v.checkRule(ifStmt, chain)
	default:
		panic("unexpected node type for else")
	}
}

func (v *visitor) checkRule(ifStmt *ast.IfStmt, chain Chain) {
	failMsg := v.rule.CheckIfElse(chain, v.args)
	if failMsg == "" {
		return
	}
	if chain.HasInitializer {
		// if statement has a := initializer, so we might need to move the assignment
		// onto its own line in case the body references it
		failMsg += " (move short variable declaration to its own line if necessary)"
	}
	v.failures = append(v.failures, lint.Failure{
		Confidence: 1,
		Node:       v.target.node(ifStmt),
		Failure:    failMsg,
	})
}
