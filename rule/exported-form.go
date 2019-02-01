package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/mgechev/revive/lint"
)

// ExportedFormRule lints given else constructs.
type ExportedFormRule struct{}

// Apply applies the rule to given file.
func (r *ExportedFormRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if isTest(file) {
		return failures
	}

	fileAst := file.AST
	walker := lintExportedForm{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
		genDeclMissingComments: make(map[*ast.GenDecl]bool),
	}

	ast.Walk(&walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ExportedFormRule) Name() string {
	return "exported-form"
}

type lintExportedForm struct {
	file                   *lint.File
	fileAst                *ast.File
	lastGen                *ast.GenDecl
	genDeclMissingComments map[*ast.GenDecl]bool
	onFailure              func(lint.Failure)
}

func (w *lintExportedForm) lintFuncDoc(fn *ast.FuncDecl) {
	if !ast.IsExported(fn.Name.Name) {
		// func is unexported
		return
	}
	kind := "function"
	name := fn.Name.Name
	if fn.Recv != nil && len(fn.Recv.List) > 0 {
		// method
		kind = "method"
		recv := receiverType(fn)
		if !ast.IsExported(recv) {
			// receiver is unexported
			return
		}
		if commonMethods[name] {
			return
		}
		switch name {
		case "Len", "Less", "Swap":
			if w.file.Pkg.Sortable[recv] {
				return
			}
		}
		name = recv + "." + name
	}
	if fn.Doc == nil {
		return
	}
	s := fn.Doc.Text()
	prefix := fn.Name.Name + " "
	if !strings.HasPrefix(s, prefix) {
		w.onFailure(lint.Failure{
			Node:       fn.Doc,
			Confidence: 0.8,
			Category:   "comments",
			Failure:    fmt.Sprintf(`comment on exported %s %s should be of the form "%s..."`, kind, name, prefix),
		})
	}
}

func (w *lintExportedForm) lintTypeDoc(t *ast.TypeSpec, doc *ast.CommentGroup) {
	if !ast.IsExported(t.Name.Name) {
		return
	}
	if doc == nil {
		return
	}

	s := doc.Text()
	articles := [...]string{"A", "An", "The", "This"}
	for _, a := range articles {
		if t.Name.Name == a {
			continue
		}
		if strings.HasPrefix(s, a+" ") {
			s = s[len(a)+1:]
			break
		}
	}
	if !strings.HasPrefix(s, t.Name.Name+" ") {
		w.onFailure(lint.Failure{
			Node:       doc,
			Confidence: 1,
			Category:   "comments",
			Failure:    fmt.Sprintf(`comment on exported type %v should be of the form "%v ..." (with optional leading article)`, t.Name, t.Name),
		})
	}
}

func (w *lintExportedForm) lintValueSpecDoc(vs *ast.ValueSpec, gd *ast.GenDecl, genDeclMissingComments map[*ast.GenDecl]bool) {
	kind := "var"
	if gd.Tok == token.CONST {
		kind = "const"
	}

	// Only one name.
	name := vs.Names[0].Name
	if !ast.IsExported(name) {
		return
	}

	// If this GenDecl has parens and a comment, we don't check its comment form.
	if gd.Lparen.IsValid() && gd.Doc != nil {
		return
	}
	// The relevant text to check will be on either vs.Doc or gd.Doc.
	// Use vs.Doc preferentially.
	doc := vs.Doc
	if doc == nil {
		doc = gd.Doc
	}
	if doc == nil {
		return
	}
	prefix := name + " "
	if !strings.HasPrefix(doc.Text(), prefix) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       doc,
			Category:   "comments",
			Failure:    fmt.Sprintf(`comment on exported %s %s should be of the form "%s..."`, kind, name, prefix),
		})
	}
}

func (w *lintExportedForm) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.GenDecl:
		if v.Tok == token.IMPORT {
			return nil
		}
		// token.CONST, token.TYPE or token.VAR
		w.lastGen = v
		return w
	case *ast.FuncDecl:
		w.lintFuncDoc(v)
		// Don't proceed inside funcs.
		return nil
	case *ast.TypeSpec:
		// inside a GenDecl, which usually has the doc
		doc := v.Doc
		if doc == nil {
			doc = w.lastGen.Doc
		}
		w.lintTypeDoc(v, doc)
		// Don't proceed inside types.
		return nil
	case *ast.ValueSpec:
		w.lintValueSpecDoc(v, w.lastGen, w.genDeclMissingComments)
		return nil
	}
	return w
}
