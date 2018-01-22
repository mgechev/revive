package rule

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/mgechev/revive/linter"
)

// ExportedRule lints given else constructs.
type ExportedRule struct{}

// Apply applies the rule to given file.
func (r *ExportedRule) Apply(file *linter.File, arguments linter.Arguments) []linter.Failure {
	var failures []linter.Failure

	if isTest(file) {
		return failures
	}

	fileAst := file.AST
	walker := lintExported{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure linter.Failure) {
			failures = append(failures, failure)
		},
		genDeclMissingComments: make(map[*ast.GenDecl]bool),
	}

	ast.Walk(&walker, fileAst)

	return failures
}

// Name returns the rule name.
func (r *ExportedRule) Name() string {
	return "imports"
}

type lintExported struct {
	file                   *linter.File
	fileAst                *ast.File
	lastGen                *ast.GenDecl
	genDeclMissingComments map[*ast.GenDecl]bool
	onFailure              func(linter.Failure)
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
		// TODO: RE-ENABLE
		// switch name {
		// case "Len", "Less", "Swap":
		// 	if f.pkg.sortable[recv] {
		// 		return
		// 	}
		// }
		name = recv + "." + name
	}
	if fn.Doc == nil {
		w.onFailure(linter.Failure{
			Node:       fn,
			Failure:    fmt.Sprintf("exported %s %s should have comment or be unexported", kind, name),
			Confidence: 1,
		})
		return
	}
	s := fn.Doc.Text()
	prefix := fn.Name.Name + " "
	if !strings.HasPrefix(s, prefix) {
		w.onFailure(linter.Failure{
			Node:       fn.Doc,
			Failure:    fmt.Sprintf(`comment on exported %s %s should be of the form "%s..."`, kind, name, prefix),
			Confidence: 1,
		})
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
		w.onFailure(linter.Failure{
			Node:       id,
			Failure:    fmt.Sprintf("%s name will be used as %s.%s by other packages, and that stutters; consider calling this %s", thing, pkg, name, rem),
			Confidence: 0.8,
		})
	}
}

func (w *lintExported) lintTypeDoc(t *ast.TypeSpec, doc *ast.CommentGroup) {
	if !ast.IsExported(t.Name.Name) {
		return
	}
	if doc == nil {
		w.onFailure(linter.Failure{
			Node:       t,
			Failure:    fmt.Sprintf("exported type %v should have comment or be unexported", t.Name),
			Confidence: 1,
		})
		return
	}

	s := doc.Text()
	articles := [...]string{"A", "An", "The"}
	for _, a := range articles {
		if strings.HasPrefix(s, a+" ") {
			s = s[len(a)+1:]
			break
		}
	}
	if !strings.HasPrefix(s, t.Name.Name+" ") {
		w.onFailure(linter.Failure{
			Node:       doc,
			Failure:    fmt.Sprintf(`comment on exported type %v should be of the form "%v ..." (with optional leading article)`, t.Name, t.Name),
			Confidence: 1,
		})
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
				w.onFailure(linter.Failure{
					Node:       vs,
					Failure:    fmt.Sprintf("exported %s %s should have its own declaration", kind, n.Name),
					Confidence: 1,
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
		w.onFailure(linter.Failure{
			Node:       vs,
			Failure:    fmt.Sprintf("exported %s %s should have comment%s or be unexported", kind, name, block),
			Confidence: 1,
		})
		genDeclMissingComments[gd] = true
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
	prefix := name + " "
	if !strings.HasPrefix(doc.Text(), prefix) {
		w.onFailure(linter.Failure{
			Node:       doc,
			Failure:    fmt.Sprintf(`comment on exported %s %s should be of the form "%s..."`, kind, name, prefix),
			Confidence: 1,
		})
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
