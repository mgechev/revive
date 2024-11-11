package fixtures

import (
	"fmt"
	ast "go/ast"
	"go/token"
)

func funLengthA() (a int) {
	println()
	println()
	println()
	println()
	println()
}

func funLengthB(file *ast.File, fset *token.FileSet, lineLimit, stmtLimit int) []Message { // MATCH /maximum number of lines per function exceeded; max 5 but got 23/
	if true {
		a = b
		if false {
			c = d
			for _, f := range list {
				_, ok := f.(int64)
				if !ok {
					continue
				}
			}
			switch a {
			case 1:
				println()
			case 2:
				println()
				println()
			default:
				println()

			}
		}
	}
	return
}

func funLengthC(b []ast.Stmt) int { // MATCH /maximum number of lines per function exceeded; max 5 but got 23/
	count := 0
	for _, s := range b {
		switch stmt := s.(type) {
		case *ast.BlockStmt:
			count += w.countStmts(stmt.List)
		case *ast.ForStmt, *ast.RangeStmt, *ast.IfStmt,
			*ast.SwitchStmt, *ast.TypeSwitchStmt, *ast.SelectStmt:
			count += 1 + w.countBodyListStmts(stmt)
		case *ast.CaseClause:
			count += w.countStmts(stmt.Body)
		case *ast.AssignStmt:
			count += 1 + w.countFuncLitStmts(stmt.Rhs[0])
		case *ast.GoStmt:
			count += 1 + w.countFuncLitStmts(stmt.Call.Fun)
		case *ast.DeferStmt:
			count += 1 + w.countFuncLitStmts(stmt.Call.Fun)
		default:
			fmt.Printf("%T %v\n", stmt, stmt)
			count++
		}
	}

	return count
}

func funLengthD(b []ast.Stmt) int {
	defer func() { println() }()
}

func funLengthE(b []ast.Stmt) int { // MATCH /maximum number of lines per function exceeded; max 5 but got 7/
	defer func() {
		if true {
			println()
		} else {
			print()
		}
	}()
}

func funLengthF(b []ast.Stmt) int {
	if true {
		println()
	} else {
		print()
	}
}

func funLengthG(b []ast.Stmt) int { // MATCH /maximum number of lines per function exceeded; max 5 but got 7/
	go func() {
		if true {
			println()
		} else {

		}
	}()
}

func funLengthH(b []ast.Stmt) int {
	go func() {}()
	println()
}
