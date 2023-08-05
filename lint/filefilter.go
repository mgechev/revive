package lint

import (
	"fmt"
	"regexp"
	"strings"
)

// FileFilter - file filter to exclude some files for rule
// supports whole
// 1. file/dir names : pkg/mypkg/my.go,
// 2. globs: **/*.pb.go,
// 3. regexes (~ prefix) ~-tmp\.\d+\.go
// 4. special test marker `TEST` - treats as `~_test\.go`
type FileFilter struct {
	// raw definition of filter inside config
	raw string
	// don't care what was at start, will use regexes inside
	rx *regexp.Regexp
	// marks that it was empty rule that matches everything
	matchesAll bool
}

// ParseFileFilter - creates [FileFilter] for given raw filter
// if empty string, or `*`, or `~` is used it means "always true"
// while regexp could be invalid, it could return it's compilation error
func ParseFileFilter(cfgFilter string) (*FileFilter, error) {
	cfgFilter = strings.TrimSpace(cfgFilter)
	result := new(FileFilter)
	result.raw = cfgFilter
	result.matchesAll = len(result.raw) == 0 || result.raw == "*" || result.raw == "~"
	if !result.matchesAll {
		if err := result.prepareRegexp(); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (ff *FileFilter) String() string { return ff.raw }

// MatchFileName - checks if file name matches filter
func (ff *FileFilter) MatchFileName(name string) bool {
	if ff.matchesAll {
		return true
	}
	name = strings.ReplaceAll(name, "\\", "/")
	return ff.rx.MatchString(name)
}

// Match - checks if given [File] matches filter
func (ff *FileFilter) Match(f *File) bool {
	return ff.MatchFileName(f.Name)
}

var fileFilterInvalidGlobRegexp = regexp.MustCompile(`[^/]\*\*[^/]`)
var escapeRegexSymbols = ".+{}()[]^$"

func (ff *FileFilter) prepareRegexp() error {
	var err error
	var src = ff.raw
	if src == "TEST" {
		src = "~_test\\.go"
	}
	if strings.HasPrefix(src, "~") {
		ff.rx, err = regexp.Compile(src[1:])
		if err != nil {
			return fmt.Errorf("invalid file filter [%s], regexp compile error: [%v]", ff.raw, err)
		}
		return nil
	}
	/* globs */
	if strings.Contains(src, "*") {
		if fileFilterInvalidGlobRegexp.MatchString(src) {
			return fmt.Errorf("invalid file filter [%s], invalid glob pattern", ff.raw)
		}
		var rxBuild strings.Builder
		rxBuild.WriteByte('^')
		wasStar := false
		justDirGlob := false
		for _, c := range src {
			if c == '*' {
				if wasStar {
					rxBuild.WriteString(`[\s\S]*`)
					wasStar = false
					justDirGlob = true
					continue
				}
				wasStar = true
				continue
			}
			if wasStar {
				rxBuild.WriteString("[^/]*")
				wasStar = false
			}
			if strings.ContainsRune(escapeRegexSymbols, c) {
				rxBuild.WriteByte('\\')
			}
			rxBuild.WriteRune(c)
			if c == '/' && justDirGlob {
				rxBuild.WriteRune('?')
			}
			justDirGlob = false
		}
		if wasStar {
			rxBuild.WriteString("[^/]*")
		}
		rxBuild.WriteByte('$')
		ff.rx, err = regexp.Compile(rxBuild.String())
		if err != nil {
			return fmt.Errorf("invalid file filter [%s], regexp compile error after glob expand: [%v]", ff.raw, err)
		}
		return nil
	}

	// it's whole file mask, just escape dots and normilze separators
	fillRx := src
	fillRx = strings.ReplaceAll(fillRx, "\\", "/")
	fillRx = strings.ReplaceAll(fillRx, ".", `\.`)
	fillRx = "^" + fillRx + "$"
	ff.rx, err = regexp.Compile(fillRx)
	if err != nil {
		return fmt.Errorf("invalid file filter [%s], regexp compile full path: [%v]", ff.raw, err)
	}
	return nil
}
