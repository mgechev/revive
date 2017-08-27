package visitors

import (
	"fmt"
	"go/ast"
)

// SyntaxVisitor implements a visitor which knows how to handle the individual
// Go lang syntax constructs.
type SyntaxVisitor struct {
	Impl Visitor
}

// Visit accepts an ast.Node and traverse its children.
func (w *SyntaxVisitor) Visit(node ast.Node) {
	if node == nil {
		return
	}
	switch x := node.(type) {
	case *ast.UnaryExpr:
		w.Impl.VisitUnaryExpr(x)
		break
	case *ast.StructType:
		w.Impl.VisitStructType(x)
		break
	case *ast.File:
		w.Impl.VisitFile(x)
		break
	case *ast.SendStmt:
		w.Impl.VisitSendStmt(x)
		break
	case *ast.ReturnStmt:
		w.Impl.VisitReturnStmt(x)
		break
	case *ast.ArrayType:
		w.Impl.VisitArrayType(x)
		break
	case *ast.AssignStmt:
		w.Impl.VisitAssignStmt(x)
		break
	case *ast.BasicLit:
		w.Impl.VisitBasicLit(x)
		break
	case *ast.BinaryExpr:
		w.Impl.VisitBinaryExpr(x)
		break
	case *ast.BlockStmt:
		w.Impl.VisitBlockStmt(x)
		break
	case *ast.BranchStmt:
		w.Impl.VisitBranchStmt(x)
		break
	case *ast.CallExpr:
		w.Impl.VisitCallExpr(x)
		break
	case *ast.CaseClause:
		w.Impl.VisitCaseClause(x)
		break
	case *ast.ChanType:
		w.Impl.VisitChanType(x)
		break
	case *ast.CommClause:
		w.Impl.VisitCommClause(x)
		break
	case *ast.Comment:
		w.Impl.VisitComment(x)
		break
	case *ast.CommentGroup:
		w.Impl.VisitCommentGroup(x)
		break
	case *ast.CompositeLit:
		w.Impl.VisitCompositeLit(x)
		break
	case *ast.DeclStmt:
		w.Impl.VisitDeclStmt(x)
		break
	case *ast.DeferStmt:
		w.Impl.VisitDeferStmt(x)
		break
	case *ast.Ellipsis:
		w.Impl.VisitEllipsis(x)
		break
	case *ast.EmptyStmt:
		w.Impl.VisitEmptyStmt(x)
		break
	case *ast.ExprStmt:
		w.Impl.VisitExprStmt(x)
		break
	case *ast.Field:
		w.Impl.VisitField(x)
		break
	case *ast.FieldList:
		w.Impl.VisitFieldList(x)
		break
	case *ast.ForStmt:
		w.Impl.VisitForStmt(x)
		break
	case *ast.FuncDecl:
		w.Impl.VisitFuncDecl(x)
		break
	case *ast.FuncLit:
		w.Impl.VisitFuncLit(x)
		break
	case *ast.FuncType:
		w.Impl.VisitFuncType(x)
		break
	case *ast.GenDecl:
		w.Impl.VisitGenDecl(x)
		break
	case *ast.GoStmt:
		w.Impl.VisitGoStmt(x)
		break
	case *ast.Ident:
		w.Impl.VisitIdent(x)
		break
	case *ast.IfStmt:
		w.Impl.VisitIfStmt(x)
		break
	case *ast.ImportSpec:
		w.Impl.VisitImportSpec(x)
		break
	case *ast.IncDecStmt:
		w.Impl.VisitIncDecStmt(x)
		break
	case *ast.IndexExpr:
		w.Impl.VisitIndexExpr(x)
		break
	case *ast.InterfaceType:
		w.Impl.VisitInterfaceType(x)
		break
	case *ast.KeyValueExpr:
		w.Impl.VisitKeyValueExpr(x)
		break
	case *ast.LabeledStmt:
		w.Impl.VisitLabeledStmt(x)
		break
	case *ast.MapType:
		w.Impl.VisitMapType(x)
		break
	case *ast.ParenExpr:
		w.Impl.VisitParenExpr(x)
		break
	case *ast.SelectorExpr:
		w.Impl.VisitSelectorExpr(x)
		break
	case *ast.SliceExpr:
		w.Impl.VisitSliceExpr(x)
		break
	case *ast.SwitchStmt:
		w.Impl.VisitSwitchStmt(x)
		break
	case *ast.TypeAssertExpr:
		w.Impl.VisitTypeAssertExpr(x)
		break
	case *ast.TypeSwitchStmt:
		w.Impl.VisitTypeSwitchStmt(x)
		break
	case *ast.StarExpr:
		w.Impl.VisitStarExpr(x)
		break
	case *ast.SelectStmt:
		w.Impl.VisitSelectStmt(x)
		break
	case *ast.RangeStmt:
		w.Impl.VisitRangeStmt(x)
		break
	case *ast.ValueSpec:
		w.Impl.VisitValueSpec(x)
		break
	case *ast.TypeSpec:
		w.Impl.VisitTypeSpec(x)
		break
	case *ast.Package:
		w.Impl.VisitPackage(x)
		break
	case *ast.BadStmt:
		w.Impl.VisitBadStmt(x)
		break
	case *ast.BadDecl:
		w.Impl.VisitBadDecl(x)
		break
	default:
		panic(fmt.Sprintf("ast.Walk: unexpected node type %T", node))
	}
}

// VisitUnaryExpr visits an unary expression.
func (w *SyntaxVisitor) VisitUnaryExpr(node *ast.UnaryExpr) {
	w.Impl.Visit(node.X)
}

func (w *SyntaxVisitor) VisitStructType(node *ast.StructType) {
	w.Impl.Visit(node.Fields)
}

func (w *SyntaxVisitor) VisitSendStmt(node *ast.SendStmt) {
	w.Impl.Visit(node.Chan)
	w.Impl.Visit(node.Value)
}

func (w *SyntaxVisitor) VisitReturnStmt(node *ast.ReturnStmt) {
	for _, x := range node.Results {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitArrayType(node *ast.ArrayType) {
	if node.Len != nil {
		w.Impl.Visit(node.Len)
	}
	w.Impl.Visit(node.Elt)
}

func (w *SyntaxVisitor) VisitAssignStmt(node *ast.AssignStmt) {
	for _, x := range node.Lhs {
		w.Impl.Visit(x)
	}
	for _, x := range node.Rhs {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitBasicLit(node *ast.BasicLit) {
}

func (w *SyntaxVisitor) VisitBinaryExpr(node *ast.BinaryExpr) {
	w.Impl.Visit(node.X)
	w.Impl.Visit(node.Y)
}

func (w *SyntaxVisitor) VisitBlockStmt(node *ast.BlockStmt) {
	for _, x := range node.List {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitFile(node *ast.File) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	w.Impl.Visit(node.Name)

	for _, x := range node.Decls {
		w.Impl.Visit(x)
	}
	// don't walk n.Comments - they have been
	// visited already through the individual
	// nodes
}

func (w *SyntaxVisitor) VisitDeclStmt(node *ast.DeclStmt) {
	w.Impl.Visit(node.Decl)
}

func (w *SyntaxVisitor) VisitDecl(node *ast.Decl) {
}

func (w *SyntaxVisitor) VisitBranchStmt(node *ast.BranchStmt) {
	if node.Label != nil {
		w.Impl.Visit(node.Label)
	}
}

func (w *SyntaxVisitor) VisitCallExpr(node *ast.CallExpr) {
	w.Impl.Visit(node.Fun)
	for _, x := range node.Args {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitCaseClause(node *ast.CaseClause) {
	for _, x := range node.List {
		w.Impl.Visit(x)
	}

	for _, x := range node.Body {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitChanType(node *ast.ChanType) {
	w.Impl.Visit(node.Value)
}

func (w *SyntaxVisitor) VisitCommClause(node *ast.CommClause) {
	if node.Comm != nil {
		w.Impl.Visit(node.Comm)
	}

	for _, x := range node.Body {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitComment(node *ast.Comment) {
}

func (w *SyntaxVisitor) VisitCommentGroup(node *ast.CommentGroup) {
	for _, x := range node.List {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitCompositeLit(node *ast.CompositeLit) {
	if node.Type != nil {
		w.Impl.Visit(node.Type)
	}

	for _, x := range node.Elts {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitDeferStmt(node *ast.DeferStmt) {
	w.Impl.Visit(node.Call)
}

func (w *SyntaxVisitor) VisitEllipsis(node *ast.Ellipsis) {
	w.Impl.Visit(node.Elt)
}

func (w *SyntaxVisitor) VisitEmptyStmt(node *ast.EmptyStmt) {
}

func (w *SyntaxVisitor) VisitExprStmt(node *ast.ExprStmt) {
	w.Impl.Visit(node.X)
}

func (w *SyntaxVisitor) VisitField(node *ast.Field) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	for _, x := range node.Names {
		w.Impl.Visit(x)
	}
	w.Impl.Visit(node.Type)
	if node.Tag != nil {
		w.Impl.Visit(node.Tag)
	}
	if node.Comment != nil {
		w.Impl.Visit(node.Comment)
	}
}

func (w *SyntaxVisitor) VisitFieldList(node *ast.FieldList) {
	for _, x := range node.List {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitForStmt(node *ast.ForStmt) {
	if node.Init != nil {
		w.Impl.Visit(node.Init)
	}
	if node.Cond != nil {
		w.Impl.Visit(node.Cond)
	}
	if node.Post != nil {
		w.Impl.Visit(node.Post)
	}
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitFuncDecl(node *ast.FuncDecl) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	if node.Recv != nil {
		w.Impl.Visit(node.Recv)
	}
	w.Impl.Visit(node.Name)
	w.Impl.Visit(node.Type)
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitFuncLit(node *ast.FuncLit) {
	w.Impl.Visit(node.Type)
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitFuncType(node *ast.FuncType) {
	if node.Params != nil {
		w.Impl.Visit(node.Params)
	}
	if node.Results != nil {
		w.Impl.Visit(node.Results)
	}
}

func (w *SyntaxVisitor) VisitGenDecl(node *ast.GenDecl) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	for _, x := range node.Specs {
		w.Impl.Visit(x)
	}
}

func (w *SyntaxVisitor) VisitGoStmt(node *ast.GoStmt) {
	w.Impl.Visit(node.Call)
}

func (w *SyntaxVisitor) VisitIdent(node *ast.Ident) {}

func (w *SyntaxVisitor) VisitIfStmt(node *ast.IfStmt) {
	if node.Init != nil {
		w.Impl.Visit(node.Init)
	}
	w.Impl.Visit(node.Cond)
	w.Impl.Visit(node.Body)
	if node.Else != nil {
		w.Impl.Visit(node.Else)
	}
}

func (w *SyntaxVisitor) VisitImportSpec(node *ast.ImportSpec) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	if node.Name != nil {
		w.Impl.Visit(node.Name)
	}
	w.Impl.Visit(node.Path)
	if node.Comment != nil {
		w.Impl.Visit(node.Comment)
	}
}

func (w *SyntaxVisitor) VisitIncDecStmt(node *ast.IncDecStmt) {
	w.Impl.Visit(node.X)
}

func (w *SyntaxVisitor) VisitIndexExpr(node *ast.IndexExpr) {
	w.Impl.Visit(node.X)
	w.Impl.Visit(node.Index)
}

func (w *SyntaxVisitor) VisitInterfaceType(node *ast.InterfaceType) {
	w.Impl.Visit(node.Methods)
}

func (w *SyntaxVisitor) VisitKeyValueExpr(node *ast.KeyValueExpr) {
	w.Impl.Visit(node.Key)
	w.Impl.Visit(node.Value)
}

func (w *SyntaxVisitor) VisitLabeledStmt(node *ast.LabeledStmt) {
	w.Impl.Visit(node.Label)
	w.Impl.Visit(node.Stmt)
}

func (w *SyntaxVisitor) VisitMapType(node *ast.MapType) {
	w.Impl.Visit(node.Key)
	w.Impl.Visit(node.Value)
}

func (w *SyntaxVisitor) VisitParenExpr(node *ast.ParenExpr) {
	w.Impl.Visit(node.X)
}

func (w *SyntaxVisitor) VisitSelectorExpr(node *ast.SelectorExpr) {
	w.Impl.Visit(node.X)
	w.Impl.Visit(node.Sel)
}

func (w *SyntaxVisitor) VisitSliceExpr(node *ast.SliceExpr) {
	w.Impl.Visit(node.X)
	if node.Low != nil {
		w.Impl.Visit(node.Low)
	}
	if node.High != nil {
		w.Impl.Visit(node.High)
	}
	if node.Max != nil {
		w.Impl.Visit(node.Max)
	}
}

func (w *SyntaxVisitor) VisitTypeAssertExpr(node *ast.TypeAssertExpr) {
	w.Impl.Visit(node.X)
	if node.Type != nil {
		w.Impl.Visit(node.Type)
	}
}

func (w *SyntaxVisitor) VisitStarExpr(node *ast.StarExpr) {
	w.Impl.Visit(node.X)
}

func (w *SyntaxVisitor) VisitSwitchStmt(node *ast.SwitchStmt) {
	if node.Init != nil {
		w.Impl.Visit(node.Init)
	}
	if node.Tag != nil {
		w.Impl.Visit(node.Tag)
	}
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitTypeSwitchStmt(node *ast.TypeSwitchStmt) {
	if node.Init != nil {
		w.Impl.Visit(node.Init)
	}
	w.Impl.Visit(node.Assign)
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitSelectStmt(node *ast.SelectStmt) {
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitRangeStmt(node *ast.RangeStmt) {
	if node.Key != nil {
		w.Impl.Visit(node.Key)
	}
	if node.Value != nil {
		w.Impl.Visit(node.Value)
	}
	w.Impl.Visit(node.X)
	w.Impl.Visit(node.Body)
}

func (w *SyntaxVisitor) VisitValueSpec(node *ast.ValueSpec) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	for _, x := range node.Names {
		w.Impl.Visit(x)
	}
	if node.Type != nil {
		w.Impl.Visit(node.Type)
	}
	for _, x := range node.Values {
		w.Impl.Visit(x)
	}
	if node.Comment != nil {
		w.Impl.Visit(node.Comment)
	}
}

func (w *SyntaxVisitor) VisitTypeSpec(node *ast.TypeSpec) {
	if node.Doc != nil {
		w.Impl.Visit(node.Doc)
	}
	w.Impl.Visit(node.Name)
	w.Impl.Visit(node.Type)
	if node.Comment != nil {
		w.Impl.Visit(node.Comment)
	}
}

func (w *SyntaxVisitor) VisitPackage(node *ast.Package) {
	for _, f := range node.Files {
		w.Impl.Visit(f)
	}
}

func (w *SyntaxVisitor) VisitBadStmt(node *ast.BadStmt) {}

func (w *SyntaxVisitor) VisitBadDecl(node *ast.BadDecl) {}
