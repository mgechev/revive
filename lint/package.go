package lint

import (
	"errors"
	"go/ast"
	"go/importer"
	"go/token"
	"go/types"
	"sync"

	goversion "github.com/hashicorp/go-version"

	"github.com/mgechev/revive/internal/typeparams"
)

// Package represents a package in the project.
type Package struct {
	fset      *token.FileSet
	files     map[string]*File
	goVersion *goversion.Version

	typesPkg  *types.Package
	typesInfo *types.Info

	// sortable is the set of types in the package that implement sort.Interface.
	sortable map[string]bool
	// main is whether this is a "main" package.
	main int
	sync.RWMutex
}

var (
	trueValue  = 1
	falseValue = 2

	go121 = goversion.Must(goversion.NewVersion("1.21"))
	go122 = goversion.Must(goversion.NewVersion("1.22"))
)

// Files return package's files.
func (p *Package) Files() map[string]*File {
	return p.files
}

// IsMain returns if that's the main package.
func (p *Package) IsMain() bool {
	p.Lock()
	defer p.Unlock()

	if p.main == trueValue {
		return true
	} else if p.main == falseValue {
		return false
	}
	for _, f := range p.files {
		if f.isMain() {
			p.main = trueValue
			return true
		}
	}
	p.main = falseValue
	return false
}

// TypesPkg yields information on this package
func (p *Package) TypesPkg() *types.Package {
	p.RLock()
	defer p.RUnlock()
	return p.typesPkg
}

// TypesInfo yields type information of this package identifiers
func (p *Package) TypesInfo() *types.Info {
	p.RLock()
	defer p.RUnlock()
	return p.typesInfo
}

// Sortable yields a map of sortable types in this package
func (p *Package) Sortable() map[string]bool {
	p.RLock()
	defer p.RUnlock()
	return p.sortable
}

// TypeCheck performs type checking for given package.
func (p *Package) TypeCheck() error {
	p.Lock()
	defer p.Unlock()

	// If type checking has already been performed
	// skip it.
	if p.typesInfo != nil || p.typesPkg != nil {
		return nil
	}
	config := &types.Config{
		// By setting a no-op error reporter, the type checker does as much work as possible.
		Error:    func(error) {},
		Importer: importer.Default(),
	}
	info := &types.Info{
		Types:  map[ast.Expr]types.TypeAndValue{},
		Defs:   map[*ast.Ident]types.Object{},
		Uses:   map[*ast.Ident]types.Object{},
		Scopes: map[ast.Node]*types.Scope{},
	}
	var anyFile *File
	var astFiles []*ast.File
	for _, f := range p.files {
		anyFile = f
		astFiles = append(astFiles, f.AST)
	}

	if anyFile == nil {
		// this is unlikely to happen, but technically guarantees anyFile to not be nil
		return errors.New("no ast.File found")
	}

	typesPkg, err := check(config, anyFile.AST.Name.Name, p.fset, astFiles, info)

	// Remember the typechecking info, even if config.Check failed,
	// since we will get partial information.
	p.typesPkg = typesPkg
	p.typesInfo = info

	return err
}

// check function encapsulates the call to go/types.Config.Check method and
// recovers if the called method panics (see issue #59)
func check(config *types.Config, n string, fset *token.FileSet, astFiles []*ast.File, info *types.Info) (p *types.Package, err error) {
	defer func() {
		if r := recover(); r != nil {
			err, _ = r.(error)
			p = nil
			return
		}
	}()

	return config.Check(n, fset, astFiles, info)
}

// TypeOf returns the type of expression.
func (p *Package) TypeOf(expr ast.Expr) types.Type {
	if p.typesInfo == nil {
		return nil
	}
	return p.typesInfo.TypeOf(expr)
}

type walker struct {
	nmap map[string]int
	has  map[string]int
}

// bitfield for which methods exist on each type.
const (
	bfLen = 1 << iota
	bfLess
	bfSwap
)

func (w *walker) Visit(n ast.Node) ast.Visitor {
	fn, ok := n.(*ast.FuncDecl)
	if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
		return w
	}

	recvType := typeparams.ReceiverType(fn)
	bf := getBitfieldForFunction(fn)

	w.has[recvType] |= bf

	return w
}

func (p *Package) scanSortable() {
	p.sortable = map[string]bool{}

	nmap := map[string]int{"Len": bfLen, "Less": bfLess, "Swap": bfSwap}
	has := map[string]int{}
	for _, f := range p.files {
		ast.Walk(&walker{nmap: nmap, has: has}, f.AST)
	}
	for typ, ms := range has {
		if ms == bfLen|bfLess|bfSwap {
			p.sortable[typ] = true
		}
	}
}

func (p *Package) lint(rules []Rule, config Config, failures chan Failure) {
	p.scanSortable()
	var wg sync.WaitGroup
	for _, file := range p.files {
		wg.Add(1)
		go (func(file *File) {
			file.lint(rules, config, failures)
			wg.Done()
		})(file)
	}
	wg.Wait()
}

// IsAtLeastGo121 returns true if the Go version for this package is 1.21 or higher, false otherwise
func (p *Package) IsAtLeastGo121() bool {
	return p.goVersion.GreaterThanOrEqual(go121)
}

// IsAtLeastGo122 returns true if the Go version for this package is 1.22 or higher, false otherwise
func (p *Package) IsAtLeastGo122() bool {
	return p.goVersion.GreaterThanOrEqual(go122)
}

func getBitfieldForFunction(fn *ast.FuncDecl) int {
	switch {
	case funcSignatureIs(fn, "Len", []string{}, []string{"int"}):
		return bfLen
	case funcSignatureIs(fn, "Less", []string{"int", "int"}, []string{"bool"}):
		return bfLess
	case funcSignatureIs(fn, "Swap", []string{"int", "int"}, []string{}):
		return bfSwap
	default:
		return 0
	}
}

// funcSignatureIs returns true if the given func decl satisfies has a signature characterized
// by the given name, parameters types and return types; false otherwise
func funcSignatureIs(funcDecl *ast.FuncDecl, wantName string, wantParametersTypes, wantResultsTypes []string) bool {
	if wantName != funcDecl.Name.String() {
		return false // func name doesn't match expected one
	}

	funcParametersTypes := getTypeNames(funcDecl.Type.Params)
	if len(wantParametersTypes) != len(funcParametersTypes) {
		return false // func has not the expected number of parameters
	}

	funcResultsTypes := getTypeNames(funcDecl.Type.Results)
	if len(wantResultsTypes) != len(funcResultsTypes) {
		return false // func has not the expected number of return values
	}

	for i, wantType := range wantParametersTypes {
		if wantType != funcParametersTypes[i] {
			return false // type of a func's parameter does not match the type of the corresponding expected parameter
		}
	}

	for i, wantType := range wantResultsTypes {
		if wantType != funcResultsTypes[i] {
			return false // type of a func's return value does not match the type of the corresponding expected return value
		}
	}

	return true
}

func getTypeNames(fields *ast.FieldList) []string {
	result := []string{}

	if fields == nil {
		return result
	}

	for _, field := range fields.List {
		typeName := field.Type.(*ast.Ident).Name
		if field.Names == nil { // unnamed field
			result = append(result, typeName)
			continue
		}

		for range field.Names { // add one type name for each field name
			result = append(result, typeName)
		}
	}

	return result
}
