// Package revivelib provides revive's linting functionality as a lib.
package revivelib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/mgechev/dots"
	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/logging"
)

// Revive is responsible for running linters and formatters
// and returning a set of results.
type Revive struct {
	config       *lint.Config
	lintingRules []lint.Rule
	logger       *log.Logger
	maxOpenFiles int
}

// New creates a new instance of Revive lint runner.
func New(
	conf *lint.Config,
	setExitStatus bool,
	maxOpenFiles int,
	extraRules ...ExtraRule,
) (*Revive, error) {
	logger, err := logging.GetLogger()
	if err != nil {
		return nil, fmt.Errorf("initializing revive - getting logger: %w", err)
	}

	if setExitStatus {
		conf.ErrorCode = 1
		conf.WarningCode = 1
	}

	extraRuleInstances := make([]lint.Rule, len(extraRules))
	for i, extraRule := range extraRules {
		extraRuleInstances[i] = extraRule.Rule

		ruleName := extraRule.Rule.Name()

		_, isRuleAlreadyConfigured := conf.Rules[ruleName]
		if !isRuleAlreadyConfigured {
			conf.Rules[ruleName] = extraRule.DefaultConfig
		}
	}

	lintingRules, err := config.GetLintingRules(conf, extraRuleInstances)
	if err != nil {
		return nil, fmt.Errorf("initializing revive - getting lint rules: %w", err)
	}

	logger.Println("Config loaded")

	return &Revive{
		logger:       logger,
		config:       conf,
		lintingRules: lintingRules,
		maxOpenFiles: maxOpenFiles,
	}, nil
}

// Lint the included patterns, skipping excluded ones
func (r *Revive) Lint(patterns ...*LintPattern) (<-chan lint.Failure, error) {
	includePatterns := []string{}
	excludePatterns := []string{}

	for _, lintpkg := range patterns {
		if lintpkg.IsExclude() {
			excludePatterns = append(excludePatterns, lintpkg.GetPattern())
		} else {
			includePatterns = append(includePatterns, lintpkg.GetPattern())
		}
	}

	if len(excludePatterns) == 0 { // if no excludes were set
		excludePatterns = r.config.Exclude // use those from the configuration
	}

	packages, err := getPackages(includePatterns, excludePatterns)
	if err != nil {
		return nil, fmt.Errorf("linting - getting packages: %w", err)
	}

	revive := lint.New(func(file string) ([]byte, error) {
		contents, err := os.ReadFile(file)

		if err != nil {
			return nil, fmt.Errorf("reading file %v: %w", file, err)
		}

		return contents, nil
	}, r.maxOpenFiles)

	failures, err := revive.Lint(packages, r.lintingRules, *r.config)
	if err != nil {
		return nil, fmt.Errorf("linting - retrieving failures channel: %w", err)
	}

	return failures, nil
}

// Format gets the output for a given failures channel from Lint.
func (r *Revive) Format(
	formatterName string,
	failuresChan <-chan lint.Failure,
) (string, int, error) {
	conf := r.config
	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	formatter, err := config.GetFormatter(formatterName)
	if err != nil {
		return "", 0, fmt.Errorf("formatting - getting formatter: %w", err)
	}

	var (
		output    string
		formatErr error
	)

	go func() {
		output, formatErr = formatter.Format(formatChan, *conf)

		exitChan <- true
	}()

	exitCode := 0

	for failure := range failuresChan {
		if failure.Confidence < conf.Confidence {
			continue
		}

		if exitCode == 0 {
			exitCode = conf.WarningCode
		}

		if c, ok := conf.Rules[failure.RuleName]; ok && c.Severity == lint.SeverityError {
			exitCode = conf.ErrorCode
		}

		if c, ok := conf.Directives[failure.RuleName]; ok && c.Severity == lint.SeverityError {
			exitCode = conf.ErrorCode
		}

		formatChan <- failure
	}

	close(formatChan)
	<-exitChan

	if formatErr != nil {
		return "", exitCode, fmt.Errorf("formatting: %w", err)
	}

	return output, exitCode, nil
}

func getPackages(includePatterns []string, excludePatterns ArrayFlags) ([][]string, error) {
	globs := normalizeSplit(includePatterns)
	if len(globs) == 0 {
		globs = append(globs, ".")
	}

	globs, skips, err := prepareSkips(globs, normalizeSplit(excludePatterns))
	if err != nil {
		return nil, fmt.Errorf("prepare skips - resolving excludes before dots: %w", err)
	}

	packages, err := dots.ResolvePackages(globs, skips)
	if err != nil {
		return nil, fmt.Errorf("getting packages - resolving packages in dots: %w", err)
	}

	return packages, nil
}

func prepareSkips(globs, excludes []string) ([]string, []string, error) {
	var skips []string
	for _, path := range globs {
		var basepath string
		basepath, _ = doublestar.SplitPattern(path)
		fsys := os.DirFS(basepath)
		for _, skip := range excludes {
			matches, err := doublestar.Glob(fsys, skip)
			if err != nil {
				return nil, nil, fmt.Errorf("Skips Error: %v", err)
			}
			for _, match := range matches {
				path = basepath + "/" + match
				// create skip only for .go files
				if filepath.Ext(path) == ".go" {
					skips = append(skips, path)
				}
			}
		}
	}
	return globs, skips, nil
}

func normalizeSplit(strs []string) []string {
	res := []string{}

	for _, s := range strs {
		t := strings.Trim(s, " \t")
		if len(t) > 0 {
			res = append(res, t)
		}
	}

	return res
}

// ArrayFlags type for string list.
type ArrayFlags []string

func (i *ArrayFlags) String() string {
	return strings.Join([]string(*i), " ")
}

// Set value for array flags.
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
}
