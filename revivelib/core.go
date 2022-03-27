package revivelib

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/mgechev/dots"
	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/logging"
	"github.com/pkg/errors"
)

// Revive is responsible for running linters and formatters
// and returning a set of results.
type Revive struct {
	config          *lint.Config
	lintingRules    []lint.Rule
	logger          *log.Logger
	maxOpenFiles    int
	excludePatterns ArrayFlags
}

// New creates a new instance of Revive lint runner.
func New(
	conf *lint.Config,
	setExitStatus bool,
	maxOpenFiles int,
	excludePatterns ArrayFlags,
	extraRules ...ExtraRule,
) (*Revive, error) {
	log, err := logging.GetLogger()
	if err != nil {
		return nil, errors.Wrap(err, "initializing revive - getting logger")
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
		return nil, errors.Wrap(err, "initializing revive - gettint lint rules")
	}

	log.Println("Config loaded")

	if len(excludePatterns) == 0 { // if no excludes were set in the command line
		excludePatterns = conf.Exclude // use those from the configuration
	}

	return &Revive{
		logger:          log,
		config:          conf,
		lintingRules:    lintingRules,
		maxOpenFiles:    maxOpenFiles,
		excludePatterns: excludePatterns,
	}, nil
}

// Lint files in the specified paths.
func (r *Revive) Lint(includePatterns ...string) (<-chan lint.Failure, error) {
	packages, err := getPackages(includePatterns, r.excludePatterns)
	if err != nil {
		return nil, errors.Wrap(err, "linting - getting packages")
	}

	revive := lint.New(func(file string) ([]byte, error) {
		contents, err := ioutil.ReadFile(file)

		if err != nil {
			return nil, errors.Wrap(err, "reading file "+file)
		}

		return contents, nil
	}, r.maxOpenFiles)

	failures, err := revive.Lint(packages, r.lintingRules, *r.config)
	if err != nil {
		return nil, errors.Wrap(err, "linting - retrieving failures channel")
	}

	return failures, nil
}

// GetLintFailures gets the list of failures for a given failures channel from Lint.
func (r *Revive) GetLintFailures(failuresChan <-chan lint.Failure) []lint.Failure {
	conf := r.config

	result := []lint.Failure{}

	for failure := range failuresChan {
		if failure.Confidence < conf.Confidence {
			continue
		}

		result = append(result, failure)
	}

	return result
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
		return "", 0, errors.Wrap(err, "formatting - getting formatter")
	}

	var output string

	go func() {
		output, err = formatter.Format(formatChan, *conf)

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

	if err != nil {
		return "", exitCode, errors.Wrap(err, "formatting")
	}

	return output, exitCode, nil
}

func getPackages(includePatterns []string, excludePatterns ArrayFlags) ([][]string, error) {
	globs := normalizeSplit(includePatterns)
	if len(globs) == 0 {
		globs = append(globs, ".")
	}

	packages, err := dots.ResolvePackages(globs, normalizeSplit(excludePatterns))
	if err != nil {
		return nil, errors.Wrap(err, "getting packages - resolving packages in dots")
	}

	return packages, nil
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
