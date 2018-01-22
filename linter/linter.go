package linter

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"math"
	"regexp"
	"strings"
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

// isGenerated reports whether the source file is generated code
// according the rules from https://golang.org/s/generatedcode.
// This is inherited from the original linter.
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

// Lint lints a set of files with the specified rule.
func (l *Linter) Lint(filenames []string, ruleSet []Rule, rulesConfig RulesConfig) ([]Failure, error) {
	var failures []Failure
	ruleNames := []string{}
	for _, r := range ruleSet {
		ruleNames = append(ruleNames, r.Name())
	}
	pkg := &Package{
		Fset:  token.NewFileSet(),
		Files: map[string]*File{},
	}
	var pkgName string
	for _, filename := range filenames {
		content, err := l.reader(filename)
		if err != nil {
			return nil, err
		}
		if isGenerated(content) {
			continue
		}

		file, err := NewFile(filename, content, pkg)
		if err != nil {
			return nil, err
		}

		if pkgName == "" {
			pkgName = file.GetAST().Name.Name
		} else if file.GetAST().Name.Name != pkgName {
			return nil, fmt.Errorf("%s is in package %s, not %s", filename, file.GetAST().Name.Name, pkgName)
		}

		pkg.Files[filename] = file
		disabledIntervals := l.disabledIntervals(file, ruleNames)

		pkg.TypeCheck()

		for _, currentRule := range ruleSet {
			config := rulesConfig[currentRule.Name()]
			currentFailures := currentRule.Apply(file, config)
			for idx, failure := range currentFailures {
				if failure.RuleName == "" {
					failure.RuleName = currentRule.Name()
				}
				if failure.Node != nil {
					failure.Position = ToFailurePosition(failure.Node.Pos(), failure.Node.End(), file)
				}
				currentFailures[idx] = failure
			}
			currentFailures = l.filterFailures(currentFailures, disabledIntervals)
			failures = append(failures, currentFailures...)
		}
	}

	return failures, nil
}

type enableDisableConfig struct {
	enabled  bool
	position int
}

func (l *Linter) disabledIntervals(file *File, allRuleNames []string) disabledIntervalsMap {
	re := regexp.MustCompile(`^\s*revive:(enable|disable)(?:-(line|next-line))?(:|\s|$)`)

	enabledDisabledRulesMap := make(map[string][]enableDisableConfig)

	getEnabledDisabledIntervals := func() disabledIntervalsMap {
		result := make(disabledIntervalsMap)

		for ruleName, disabledArr := range enabledDisabledRulesMap {
			ruleResult := []DisabledInterval{}
			for i := 0; i < len(disabledArr); i++ {
				interval := DisabledInterval{
					RuleName: ruleName,
					From: token.Position{
						Filename: file.Name,
						Line:     disabledArr[i].position,
					},
					To: token.Position{
						Filename: file.Name,
						Line:     math.MaxInt32,
					},
				}
				if i%2 == 0 {
					ruleResult = append(ruleResult, interval)
				} else {
					ruleResult[len(ruleResult)-1].To.Line = disabledArr[i].position
				}
			}
			result[ruleName] = ruleResult
		}

		return result
	}

	handleConfig := func(isEnabled bool, line int, name string) {
		existing, ok := enabledDisabledRulesMap[name]
		if !ok {
			existing = []enableDisableConfig{}
			enabledDisabledRulesMap[name] = existing
		}
		if (len(existing) > 1 && existing[len(existing)-1].enabled == isEnabled) ||
			(len(existing) == 0 && isEnabled) {
			return
		}
		existing = append(existing, enableDisableConfig{
			enabled:  isEnabled,
			position: line,
		})
		enabledDisabledRulesMap[name] = existing
	}

	handleRules := func(filename, modifier string, isEnabled bool, line int, ruleNames []string) []DisabledInterval {
		var result []DisabledInterval
		for _, name := range ruleNames {
			if modifier == "line" {
				handleConfig(isEnabled, line, name)
				handleConfig(!isEnabled, line, name)
			} else if modifier == "next-line" {
				handleConfig(isEnabled, line+1, name)
				handleConfig(!isEnabled, line+1, name)
			} else {
				handleConfig(isEnabled, line, name)
			}
		}
		return result
	}

	handleComment := func(filename string, c *ast.CommentGroup, line int) {
		text := c.Text()
		parts := re.FindStringSubmatch(text)
		if len(parts) == 0 {
			return
		}
		str := re.FindString(text)
		ruleNamesString := strings.Split(text, str)
		ruleNames := []string{}
		if len(ruleNamesString) == 2 {
			tempNames := strings.Split(ruleNamesString[1], ",")
			for _, name := range tempNames {
				name = strings.Trim(name, "\n")
				if len(name) > 0 {
					ruleNames = append(ruleNames, name)
				}
			}
		}

		if len(ruleNames) == 0 {
			ruleNames = allRuleNames
		}

		handleRules(filename, parts[2], parts[1] == "enable", line, ruleNames)
	}

	comments := file.GetAST().Comments
	for _, c := range comments {
		handleComment(file.Name, c, file.ToPosition(c.Pos()).Line)
	}

	return getEnabledDisabledIntervals()
}

func (l *Linter) filterFailures(failures []Failure, disabledIntervals disabledIntervalsMap) []Failure {
	result := []Failure{}
	for _, failure := range failures {
		fStart := failure.Position.Start.Line
		fEnd := failure.Position.End.Line
		intervals, ok := disabledIntervals[failure.RuleName]
		if !ok {
			result = append(result, failure)
		} else {
			include := true
			for _, interval := range intervals {
				intStart := interval.From.Line
				intEnd := interval.To.Line
				if (fStart >= intStart && fStart <= intEnd) ||
					(fEnd >= intStart && fEnd <= intEnd) {
					include = false
					break
				}
			}
			if include {
				result = append(result, failure)
			}
		}
	}
	return result
}
