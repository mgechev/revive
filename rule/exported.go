package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mgechev/revive/lint"
)

// ExportedRule lints given else constructs.
type ExportedRule struct{}

// Apply applies the rule to given file.
func (r *ExportedRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	if isTest(file) {
		return failures
	}

	fileAst := file.AST
	walker := lintExported{
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
func (r *ExportedRule) Name() string {
	return "exported"
}

type lintExported struct {
	file                   *lint.File
	fileAst                *ast.File
	lastGen                *ast.GenDecl
	genDeclMissingComments map[*ast.GenDecl]bool
	onFailure              func(lint.Failure)
}

func (w *lintExported) lintFuncDoc(fn *ast.FuncDecl) {
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
		w.onFailure(lint.Failure{
			Node:       fn,
			Confidence: 1,
			Category:   "comments",
			Failure:    fmt.Sprintf("exported %s %s should have comment or be unexported", kind, name),
		})
		return
	}
}

func (w *lintExported) checkStutter(id *ast.Ident, thing string) {
	pkg, name := w.fileAst.Name.Name, id.Name
	if !ast.IsExported(name) {
		// unexported name
		return
	}
	// A name stutters if the package name is a strict prefix
	// and the next character of the name starts a new word.
	if len(name) <= len(pkg) {
		// name is too short to stutter.
		// This permits the name to be the same as the package name.
		return
	}
	if !strings.EqualFold(pkg, name[:len(pkg)]) {
		return
	}
	// We can assume the name is well-formed UTF-8.
	// If the next rune after the package name is uppercase or an underscore
	// the it's starting a new word and thus this name stutters.
	rem := name[len(pkg):]
	if next, _ := utf8.DecodeRuneInString(rem); next == '_' || unicode.IsUpper(next) {
		w.onFailure(lint.Failure{
			Node:       id,
			Confidence: 0.8,
			Category:   "naming",
			Failure:    fmt.Sprintf("%s name will be used as %s.%s by other packages, and that stutters; consider calling this %s", thing, pkg, name, rem),
		})
	}
}

func (w *lintExported) lintTypeDoc(t *ast.TypeSpec, doc *ast.CommentGroup) {
	if !ast.IsExported(t.Name.Name) {
		return
	}
	if doc == nil {
		w.onFailure(lint.Failure{
			Node:       t,
			Confidence: 1,
			Category:   "comments",
			Failure:    fmt.Sprintf("exported type %v should have comment or be unexported", t.Name),
		})
		return
	}

}

func (w *lintExported) lintValueSpecDoc(vs *ast.ValueSpec, gd *ast.GenDecl, genDeclMissingComments map[*ast.GenDecl]bool) {
	kind := "var"
	if gd.Tok == token.CONST {
		kind = "const"
	}

	if len(vs.Names) > 1 {
		// Check that none are exported except for the first.
		for _, n := range vs.Names[1:] {
			if ast.IsExported(n.Name) {
				w.onFailure(lint.Failure{
					Category:   "comments",
					Confidence: 1,
					Failure:    fmt.Sprintf("exported %s %s should have its own declaration", kind, n.Name),
					Node:       vs,
				})
				return
			}
		}
	}

	// Only one name.
	name := vs.Names[0].Name
	if !ast.IsExported(name) {
		return
	}

	if vs.Doc == nil && gd.Doc == nil {
		if genDeclMissingComments[gd] {
			return
		}
		block := ""
		if kind == "const" && gd.Lparen.IsValid() {
			block = " (or a comment on this block)"
		}
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       vs,
			Category:   "comments",
			Failure:    fmt.Sprintf("exported %s %s should have comment%s or be unexported", kind, name, block),
		})
		genDeclMissingComments[gd] = true
		return
	}
}

func (w *lintExported) Visit(n ast.Node) ast.Visitor {
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
		if v.Recv == nil {
			// Only check for stutter on functions, not methods.
			// Method names are not used package-qualified.
			w.checkStutter(v.Name, "func")
		}
		// Don't proceed inside funcs.
		return nil
	case *ast.TypeSpec:
		// inside a GenDecl, which usually has the doc
		doc := v.Doc
		if doc == nil {
			doc = w.lastGen.Doc
		}
		w.lintTypeDoc(v, doc)
		w.checkStutter(v.Name, "type")
		// Don't proceed inside types.
		return nil
	case *ast.ValueSpec:
		w.lintValueSpecDoc(v, w.lastGen, w.genDeclMissingComments)
		return nil
	}
	return w
}
