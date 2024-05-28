package lint

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	goversion "github.com/hashicorp/go-version"
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
	genHdr = []byte("// Code generated ")
	genFtr = []byte(" DO NOT EDIT.")
)

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(packages [][]string, ruleSet []Rule, config Config) (<-chan Failure, error) {
	failures := make(chan Failure)

	var wg sync.WaitGroup
	for _, pkg := range packages {
		wg.Add(1)
		go func(pkg []string) {
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

func (l *Linter) lintPackage(filenames []string, ruleSet []Rule, config Config, failures chan Failure) error {
	if len(filenames) == 0 {
		return nil
	}

	goVersion, err := detectGoVersion(filepath.Dir(filenames[0]))
	if err != nil {
		return err
	}

	pkg := &Package{
		fset:      token.NewFileSet(),
		files:     map[string]*File{},
		goVersion: goVersion,
	}
	for _, filename := range filenames {
		content, err := l.readFile(filename)
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

func detectGoVersion(dir string) (ver *goversion.Version, err error) {
	// https://github.com/golang/go/issues/44753#issuecomment-790089020
	cmd := exec.Command("go", "list", "-m", "-json")
	cmd.Dir = dir

	raw, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("command go list: %w", err)
	}

	var v struct {
		GoMod     string `json:"GoMod"`
		GoVersion string `json:"GoVersion"`
	}
	if err = json.Unmarshal(raw, &v); err != nil {
		return nil, fmt.Errorf("can't parse the output of go list: %w", err)
	}

	if v.GoMod == "" {
		// this package is outside a module, so assume
		// an old-style source directory

		if v := os.Getenv("GOVERSION"); v != "" {
			return goversion.NewVersion(strings.TrimPrefix(v, "go"))
		}

		// assume the last version that does not have generics
		return goversion.Must(goversion.NewVersion("1.17")), nil
	}

	return goversion.NewVersion(strings.TrimPrefix(v.GoVersion, "go"))
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
