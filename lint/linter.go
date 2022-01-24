package lint

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"sync"
)

// ReadFile defines an abstraction for reading files.
type ReadFile func(path string) (result []byte, err error)

type disabledIntervalsMap = map[string][]DisabledInterval

// Linter is used for linting set of files.
type Linter struct {
	reader ReadFile
}

// New creates a new Linter
func New(reader ReadFile) Linter {
	return Linter{reader: reader}
}

var (
	genHdr = []byte("// Code generated ")
	genFtr = []byte(" DO NOT EDIT.")
)

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(packages [][]string, ruleSet []Rule, config Config) (<-chan Failure, error) {
	failures := make(chan Failure)

	concurrency := runtime.NumCPU()
	jobs := make(chan []string, 1024)

	var wg sync.WaitGroup
	wg.Add(concurrency)

	for i := 1; i <= concurrency; i++ {
		go func(jobs <-chan []string) {
			for job := range jobs {
				if err := l.lintPackage(job, ruleSet, config, failures); err != nil {
					fmt.Fprintln(os.Stderr, err)
					os.Exit(1)
				}
			}
			wg.Done()
		}(jobs)
	}

	go func() {
		for _, pkg := range packages {
			jobs <- pkg
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(failures)
	}()

	return failures, nil
}

func (l *Linter) lintPackage(filenames []string, ruleSet []Rule, config Config, failures chan Failure) error {
	pkg := &Package{
		fset:  token.NewFileSet(),
		files: map[string]*File{},
		mu:    sync.Mutex{},
	}
	for _, filename := range filenames {
		content, err := l.reader(filename)
		if err != nil {
			return err
		}
		if !config.IgnoreGeneratedHeader && isGenerated(content) {
			continue
		}

		file, err := NewFile(filename, content, pkg)
		if err != nil {
			addInvalidFileFailure(filename, err.Error(), failures)
			continue
		}
		pkg.files[filename] = file
	}

	if len(pkg.files) == 0 {
		return nil
	}

	pkg.lint(ruleSet, config, failures)

	return nil
}

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
// This is inherited from the original go lint.
func isGenerated(src []byte) bool {
	sc := bufio.NewScanner(bytes.NewReader(src))
	for sc.Scan() {
		b := sc.Bytes()
		if bytes.HasPrefix(b, genHdr) && bytes.HasSuffix(b, genFtr) && len(b) >= len(genHdr)+len(genFtr) {
			return true
		}
	}
	return false
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
var errPosRegexp = regexp.MustCompile(".*:(\\d*):(\\d*):.*$")

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
		}}
}
