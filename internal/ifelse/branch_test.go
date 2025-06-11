package ifelse

import (
	"go/ast"
	"go/token"
	"testing"
)

func TestBlockBranch(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		block := &ast.BlockStmt{List: []ast.Stmt{}}
		b := BlockBranch(block)
		if b.BranchKind != Empty {
			t.Errorf("want Empty branch, got %v", b.BranchKind)
		}
	})
	t.Run("non empty", func(t *testing.T) {
		stmt := &ast.ReturnStmt{}
		block := &ast.BlockStmt{List: []ast.Stmt{stmt}}
		b := BlockBranch(block)
		if b.BranchKind != Return {
			t.Errorf("want Return branch, got %v", b.BranchKind)
		}
	})
}

func TestStmtBranch(t *testing.T) {
	cases := []struct {
		name string
		stmt ast.Stmt
		kind BranchKind
		call *Call
	}{
		{
			name: "ReturnStmt",
			stmt: &ast.ReturnStmt{},
			kind: Return,
		},
		{
			name: "BreakStmt",
			stmt: &ast.BranchStmt{Tok: token.BREAK},
			kind: Break,
		},
		{
			name: "ContinueStmt",
			stmt: &ast.BranchStmt{Tok: token.CONTINUE},
			kind: Continue,
		},
		{
			name: "GotoStmt",
			stmt: &ast.BranchStmt{Tok: token.GOTO},
			kind: Goto,
		},
		{
			name: "EmptyStmt",
			stmt: &ast.EmptyStmt{},
			kind: Empty,
		},
		{
			name: "ExprStmt with DeviatingFunc (panic)",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "panic"},
				},
			},
			kind: Panic,
			call: &Call{Name: "panic"},
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
			kind: Exit,
			call: &Call{Pkg: "os", Name: "Exit"},
		},
		{
			name: "ExprStmt with non-deviating func",
			stmt: &ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: &ast.Ident{Name: "foo"},
				},
			},
			kind: Regular,
		},
		{
			name: "LabeledStmt wrapping ReturnStmt",
			stmt: &ast.LabeledStmt{
				Label: &ast.Ident{Name: "lbl"},
				Stmt:  &ast.ReturnStmt{},
			},
			kind: Return,
		},
		{
			name: "LabeledStmt wrapping ExprStmt",
			stmt: &ast.LabeledStmt{
				Label: &ast.Ident{Name: "lbl"},
				Stmt:  &ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: "foo"}}},
			},
			kind: Regular,
		},
		{
			name: "BlockStmt with ReturnStmt",
			stmt: &ast.BlockStmt{List: []ast.Stmt{&ast.ReturnStmt{}}},
			kind: Return,
		},
		{
			name: "BlockStmt with ExprStmt",
			stmt: &ast.BlockStmt{List: []ast.Stmt{&ast.ExprStmt{X: &ast.CallExpr{Fun: &ast.Ident{Name: "foo"}}}}},
			kind: Regular,
		},
	}
	for _, c := range cases {
		b := StmtBranch(c.stmt)
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
		branch   Branch
		wantStr  string
		wantLong string
	}{
		{
			name:     "Return branch",
			branch:   Branch{BranchKind: Return},
			wantStr:  "{ ... return }",
			wantLong: "a return statement",
		},
		{
			name:     "Panic branch with Call",
			branch:   Branch{BranchKind: Panic, Call: Call{Name: "panic"}},
			wantStr:  "{ ... panic() }",
			wantLong: "call to panic function",
		},
		{
			name:     "Exit branch with Call",
			branch:   Branch{BranchKind: Exit, Call: Call{Pkg: "os", Name: "Exit"}},
			wantStr:  "{ ... os.Exit() }",
			wantLong: "call to os.Exit function",
		},
		{
			name:     "Empty branch",
			branch:   Branch{BranchKind: Empty},
			wantStr:  "{ }",
			wantLong: "an empty block",
		},
		{
			name:     "Regular branch",
			branch:   Branch{BranchKind: Regular},
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
		block []ast.Stmt
		want  bool
	}{
		{
			name:  "DeclStmt",
			block: []ast.Stmt{&ast.DeclStmt{}},
			want:  true,
		},
		{
			name:  "AssignStmt with :=",
			block: []ast.Stmt{&ast.AssignStmt{Tok: token.DEFINE}},
			want:  true,
		},
		{
			name:  "ExprStmt",
			block: []ast.Stmt{&ast.ExprStmt{}},
			want:  false,
		},
	}
	for _, tt := range tests {
		b := Branch{block: tt.block}
		if got := b.HasDecls(); got != tt.want {
			t.Errorf("%s: want HasDecls to be %v, got %v", tt.name, tt.want, got)
		}
	}
}

func TestBranch_IsShort(t *testing.T) {
	tests := []struct {
		name  string
		block []ast.Stmt
		want  bool
	}{
		{
			name:  "nil block",
			block: nil,
			want:  true,
		},
		{
			name:  "single ExprStmt",
			block: []ast.Stmt{&ast.ExprStmt{}},
			want:  true,
		},
		{
			name:  "single BlockStmt",
			block: []ast.Stmt{&ast.BlockStmt{}},
			want:  false,
		},
		{
			name:  "two short statements",
			block: []ast.Stmt{&ast.ExprStmt{}, &ast.ExprStmt{}},
			want:  true,
		},
		{
			name:  "second non-short statement",
			block: []ast.Stmt{&ast.ExprStmt{}, &ast.BlockStmt{}},
			want:  false,
		},
		{
			name:  "three statements (should return false)",
			block: []ast.Stmt{&ast.ExprStmt{}, &ast.ExprStmt{}, &ast.ExprStmt{}},
			want:  false,
		},
	}
	for _, tt := range tests {
		b := Branch{block: tt.block}
		if got := b.IsShort(); got != tt.want {
			t.Errorf("%s: want IsShort to be %v, got %v", tt.name, tt.want, got)
		}
	}
}
