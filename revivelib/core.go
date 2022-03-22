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
	Formatter    lint.Formatter
	Config       *lint.Config
	LintingRules []lint.Rule
	Logger       *log.Logger
	MaxOpenFiles int
	ExcludePaths ArrayFlags
}

// New creates a new instance of Revive lint runner.
func New(
	formatterName string,
	conf *lint.Config,
	setExitStatus bool,
	maxOpenFiles int,
	excludePaths ArrayFlags,
	extraRules ...ExtraRule,
) (*Revive, error) {
	log, err := logging.GetLogger()
	if err != nil {
		return nil, errors.Wrap(err, "initializing revive - getting logger")
	}

	formatter, err := config.GetFormatter(formatterName)
	if err != nil {
		return nil, errors.Wrap(err, "initializing revive - getting formatter")
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

	if len(excludePaths) == 0 { // if no excludes were set in the command line
		excludePaths = conf.Exclude // use those from the configuration
	}

	return &Revive{
		Logger:       log,
		Config:       conf,
		Formatter:    formatter,
		LintingRules: lintingRules,
		MaxOpenFiles: maxOpenFiles,
		ExcludePaths: excludePaths,
	}, nil
}

// Lint files in the specified paths.
func (r *Revive) Lint(files ...string) (<-chan lint.Failure, error) {
	packages, err := getPackages(files, r.ExcludePaths)
	if err != nil {
		return nil, errors.Wrap(err, "linting - getting packages")
	}

	revive := lint.New(func(file string) ([]byte, error) {
		contents, err := ioutil.ReadFile(file)

		if err != nil {
			return nil, errors.Wrap(err, "reading file "+file)
		}

		return contents, nil
	}, r.MaxOpenFiles)

	failures, err := revive.Lint(packages, r.LintingRules, *r.Config)
	if err != nil {
		return nil, errors.Wrap(err, "linting - retrieving failures channel")
	}

	return failures, nil
}

// GetLintFailures gets the list of failures for a given failures channel from Lint.
func (r *Revive) GetLintFailures(failuresChan <-chan lint.Failure) []lint.Failure {
	conf := r.Config

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
func (r *Revive) Format(failuresChan <-chan lint.Failure) (string, int, error) {
	conf := r.Config
	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	var (
		output string
		err    error
	)

	go func() {
		output, err = r.Formatter.Format(formatChan, *conf)

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

func getPackages(files []string, excludePaths ArrayFlags) ([][]string, error) {
	globs := normalizeSplit(files)
	if len(globs) == 0 {
		globs = append(globs, ".")
	}

	packages, err := dots.ResolvePackages(globs, normalizeSplit(excludePaths))
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
