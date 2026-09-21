package ifelse_test

import (
	"go/ast"
	"go/token"
	"testing"

	"github.com/mgechev/revive/internal/ifelse"
)

func TestBlockBranch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{}}
		b := ifelse.BlockBranch(block)
		if b.BranchKind != ifelse.Empty {
			t.Errorf("want Empty branch, got %v", b.BranchKind)
		}
	})
	t.Run("non empty", func(t *testing.T) {
		stmt := &ast.ReturnStmt{}
		block := &ast.BlockStmt{List: []ast.Stmt{stmt}}
		b := ifelse.BlockBranch(block)
		if b.BranchKind != ifelse.Return {
			t.Errorf("want Return branch, got %v", b.BranchKind)
		}
	})
}

func TestStmtBranch(t *testing.T) {
	cases := []struct {
		name string
		stmt ast.Stmt
		kind ifelse.BranchKind
		call *ifelse.Call
	}{
		{
			name: "ReturnStmt",
			stmt: &ast.ReturnStmt{},
			kind: ifelse.Return,
		},
		{
			name: "BreakStmt",
			stmt: &ast.BranchStmt{Tok: token.BREAK},
			kind: ifelse.Break,
		},
		{
			name: "ContinueStmt",
			stmt: &ast.BranchStmt{Tok: token.CONTINUE},
			kind: ifelse.Continue,
		},
		{
			name: "GotoStmt",
			stmt: &ast.BranchStmt{Tok: token.GOTO},
			kind: ifelse.Goto,
		},
		{
			name: "EmptyStmt",
			stmt: &ast.EmptyStmt{},
			kind: ifelse.Empty,
		},
		{
			name: "ExprStmt with DeviatingFunc (panic)",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "panic"},
				},
			},
			kind: ifelse.Panic,
			call: &ifelse.Call{Name: "panic"},
		},
		{
			name: "ExprStmt with DeviatingFunc (os.Exit)",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X:   &ast.Ident{Name: "os"},
						Sel: &ast.Ident{Name: "Exit"},
					},
				},
			},
			kind: ifelse.Exit,
			call: &ifelse.Call{Pkg: "os", Name: "Exit"},
		},
		{
			name: "ExprStmt with non-deviating func",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "foo"},
				},
			},
			kind: ifelse.Regular,
		},
		{
			name: "LabeledStmt wrapping ReturnStmt",
			stmt: &ast.LabeledStmt{
				Label: &ast.Ident{Name: "lbl"},
				Stmt:  &ast.ReturnStmt{},
			},
			kind: ifelse.Return,
		},
		{
			name: "LabeledStmt wrapping ExprStmt",
			stmt: &ast.LabeledStmt{
				Label: &ast.Ident{Name: "lbl"},
				Stmt:  &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: "foo"}}},
			},
			kind: ifelse.Regular,
		},
		{
			name: "BlockStmt with ReturnStmt",
			stmt: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			kind: ifelse.Return,
		},
		{
			name: "BlockStmt with ExprStmt",
			stmt: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: "foo"}}}}},
			kind: ifelse.Regular,
		},
	}
	for _, c := range cases {
		b := ifelse.StmtBranch(c.stmt)
		if b.BranchKind != c.kind {
			t.Errorf("%s: want %v, got %v", c.name, c.kind, b.BranchKind)
		}
		if c.call != nil {
			if b.Call != *c.call {
				t.Errorf("%s: want Call %+v, got %+v", c.name, *c.call, b.Call)
			}
		}
	}
}

func TestBranch_String_LongString(t *testing.T) {
	tests := []struct {
		name     string
		branch   ifelse.Branch
		wantStr  string
		wantLong string
	}{
		{
			name:     "Return branch",
			branch:   ifelse.Branch{BranchKind: ifelse.Return},
			wantStr:  "{ ... return }",
			wantLong: "a return statement",
		},
		{
			name:     "Panic branch with Call",
			branch:   ifelse.Branch{BranchKind: ifelse.Panic, Call: ifelse.Call{Name: "panic"}},
			wantStr:  "{ ... panic() }",
			wantLong: "call to panic function",
		},
		{
			name:     "Exit branch with Call",
			branch:   ifelse.Branch{BranchKind: ifelse.Exit, Call: ifelse.Call{Pkg: "os", Name: "Exit"}},
			wantStr:  "{ ... os.Exit() }",
			wantLong: "call to os.Exit function",
		},
		{
			name:     "Empty branch",
			branch:   ifelse.Branch{BranchKind: ifelse.Empty},
			wantStr:  "{ }",
			wantLong: "an empty block",
		},
		{
			name:     "Regular branch",
			branch:   ifelse.Branch{BranchKind: ifelse.Regular},
			wantStr:  "{ ... }",
			wantLong: "a regular statement",
		},
	}
	for _, tt := range tests {
		if got := tt.branch.String(); got != tt.wantStr {
			t.Errorf("%s: String() = %q, want %q", tt.name, got, tt.wantStr)
		}
		if got := tt.branch.LongString(); got != tt.wantLong {
			t.Errorf("%s: LongString() = %q, want %q", tt.name, got, tt.wantLong)
		}
	}
}

func TestBranch_HasDecls(t *testing.T) {
	tests := []struct {
		name  string
		block *ast.BlockStmt
		want  bool
	}{
		{
			name:  "DeclStmt",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.DeclStmt{}}},
			want:  true,
		},
		{
			name:  "AssignStmt with :=",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.AssignStmt{Tok: token.DEFINE}}},
			want:  true,
		},
		{
			name:  "ExprStmt",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
			want:  false,
		},
	}
	for _, tt := range tests {
		b := ifelse.BlockBranch(tt.block)
		if got := b.HasDecls(); got != tt.want {
			t.Errorf("%s: want HasDecls to be %v, got %v", tt.name, tt.want, got)
		}
	}
}

func TestBranch_IsShort(t *testing.T) {
	tests := []struct {
		name  string
		block *ast.BlockStmt
		want  bool
	}{
		{
			name:  "nil block",
			block: &ast.BlockStmt{},
			want:  true,
		},
		{
			name:  "single ExprStmt",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}}},
			want:  true,
		},
		{
			name:  "single BlockStmt",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.BlockStmt{}}},
			want:  false,
		},
		{
			name:  "two short statements",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}, &ast.ExprStmt{}}},
			want:  true,
		},
		{
			name:  "second non-short statement",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}, &ast.BlockStmt{}}},
			want:  false,
		},
		{
			name:  "three statements (should return false)",
			block: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{}, &ast.ExprStmt{}, &ast.ExprStmt{}}},
			want:  false,
		},
	}
	for _, tt := range tests {
		b := ifelse.BlockBranch(tt.block)
		if got := b.IsShort(); got != tt.want {
			t.Errorf("%s: want IsShort to be %v, got %v", tt.name, tt.want, got)
		}
	}
}
