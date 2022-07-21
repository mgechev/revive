package lint

import (
	"fmt"
	"go/token"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

// ReadFile defines an abstraction for reading files.
type ReadFile func(path string) (result []byte, err error)

type disabledIntervalsMap = map[string][]DisabledInterval

// Linter is used for linting set of files.
type Linter struct {
	reader         ReadFile
	fileReadTokens chan struct{}
}

// New creates a new Linter
func New(reader ReadFile, maxOpenFiles int) Linter {
	var fileReadTokens chan struct{}
	if maxOpenFiles > 0 {
		fileReadTokens = make(chan struct{}, maxOpenFiles)
	}
	return Linter{
		reader:         reader,
		fileReadTokens: fileReadTokens,
	}
}

func (l Linter) readFile(path string) (result []byte, err error) {
	if l.fileReadTokens != nil {
		// "take" a token by writing to the channel.
		// It will block if no more space in the channel's buffer
		l.fileReadTokens <- struct{}{}
		defer func() {
			// "free" a token by reading from the channel
			<-l.fileReadTokens
		}()
	}

	return l.reader(path)
}

var (
	genHdr = "// Code generated "
	genFtr = " DO NOT EDIT."
)

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(pckgs []*packages.Package, ruleSet []Rule, config Config) (<-chan Failure, error) {
	failures := make(chan Failure)

	var wg sync.WaitGroup
	for _, pkg := range pckgs {
		wg.Add(1)
		go func(pkg *packages.Package) {
			if err := l.lintPackage(pkg, ruleSet, config, failures); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer wg.Done()
		}(pkg)
	}

	go func() {
		wg.Wait()
		close(failures)
	}()

	return failures, nil
}

func (l *Linter) lintPackage(goPkg *packages.Package, ruleSet []Rule, config Config, failures chan Failure) error {
	lintPkg := NewPackage(goPkg)
	for _, file := range goPkg.Syntax {
		/*
			content, err := l.readFile(filename)
			if err != nil {
				return err
			}*/

		if !config.IgnoreGeneratedHeader && isGenerated(file.Doc.Text()) {
			continue
		}
		/*
			file, err := NewFile(filename, content, pkg)
			if err != nil {
				addInvalidFileFailure(filename, err.Error(), failures)
				continue
			}
		*/
		newFile := File{
			Name: file.Name.Name,
			Pkg:  &lintPkg,
			AST:  file,
		}
		lintPkg.files[file.Name.String()] = &newFile
	}

	if len(lintPkg.files) == 0 {
		return nil
	}

	lintPkg.lint(ruleSet, config, failures)

	return nil
}

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
// This is inherited from the original go lint.
func isGenerated(fileComment string) bool {
	return strings.HasPrefix(fileComment, genHdr) && strings.HasSuffix(fileComment, genFtr) && len(fileComment) >= len(genHdr)+len(genFtr)
}

// addInvalidFileFailure adds a failure for an invalid formatted file
func addInvalidFileFailure(filename, errStr string, failures chan Failure) {
	position := getPositionInvalidFile(filename, errStr)
	failures <- Failure{
		Confidence: 1,
		Failure:    fmt.Sprintf("invalid file %s: %v", filename, errStr),
		Category:   "validity",
		Position:   position,
	}
}

// errPosRegexp matches with an NewFile error message
// i.e. :  corrupted.go:10:4: expected '}', found 'EOF
// first group matches the line and the second group, the column
var errPosRegexp = regexp.MustCompile(`.*:(\d*):(\d*):.*$`)

// getPositionInvalidFile gets the position of the error in an invalid file
func getPositionInvalidFile(filename, s string) FailurePosition {
	pos := errPosRegexp.FindStringSubmatch(s)
	if len(pos) < 3 {
		return FailurePosition{}
	}
	line, err := strconv.Atoi(pos[1])
	if err != nil {
		return FailurePosition{}
	}
	column, err := strconv.Atoi(pos[2])
	if err != nil {
		return FailurePosition{}
	}

	return FailurePosition{
		Start: token.Position{
			Filename: filename,
			Line:     line,
			Column:   column,
		},
	}
}
