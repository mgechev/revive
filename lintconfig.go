package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mgechev/dots"

	"github.com/mgechev/revive/formatter"

	"github.com/BurntSushi/toml"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

func fail(err string) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var defaultRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.NamesRule{},
	&rule.ElseRule{},
	&rule.IfReturnRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorsRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNameRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamesRule{},
	&rule.ContextKeyTypeRule{},
	&rule.ContextArgumentsRule{},
}

var allRules = append([]lint.Rule{
	&rule.ArgumentsLimitRule{},
	&rule.CyclomaticRule{},
}, defaultRules...)

var allFormatters = []lint.Formatter{
	&formatter.CLI{},
	&formatter.JSON{},
	&formatter.Default{},
}

func getFormatters() map[string]lint.Formatter {
	result := map[string]lint.Formatter{}
	for _, f := range allFormatters {
		result[f.Name()] = f
	}
	return result
}

func getLintingRules(config *lint.Config) []lint.Rule {
	rulesMap := map[string]lint.Rule{}
	for _, r := range allRules {
		rulesMap[r.Name()] = r
	}

	lintingRules := []lint.Rule{}
	for name := range config.Rules {
		rule, ok := rulesMap[name]
		if !ok {
			fail("cannot find rule: " + name)
		}
		lintingRules = append(lintingRules, rule)
	}

	return lintingRules
}

func parseConfig(path string) *lint.Config {
	config := &lint.Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fail("cannot read the config file")
	}
	_, err = toml.Decode(string(file), config)
	if err != nil {
		fail("cannot parse the config file: " + err.Error())
	}
	return config
}

func normalizeConfig(config *lint.Config) {
	if config.Confidence == 0 {
		config.Confidence = 0.8
	}
	severity := config.Severity
	if severity != "" {
		for k, v := range config.Rules {
			if v.Severity == "" {
				v.Severity = severity
			}
			config.Rules[k] = v
		}
	}
}

func getConfig() *lint.Config {
	config := defaultConfig()
	if configPath != "" {
		config = parseConfig(configPath)
	}
	normalizeConfig(config)
	return config
}

func getFormatter() lint.Formatter {
	formatters := getFormatters()
	formatter := formatters["default"]
	if formatterName != "" {
		f, ok := formatters[formatterName]
		if !ok {
			fail("unknown formatter " + formatterName)
		}
		formatter = f
	}
	return formatter
}

func defaultConfig() *lint.Config {
	defaultConfig := lint.Config{
		Confidence: 0.0,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range defaultRules {
		defaultConfig.Rules[r.Name()] = lint.RuleConfig{}
	}
	return &defaultConfig
}

func getFiles() []string {
	globs := flag.Args()
	if len(globs) == 0 {
		fail("files not specified")
	}

	files, errs := dots.Resolve(globs, strings.Split(excludePaths, " "))
	if errs != nil && len(errs) > 0 {
		err := ""
		for _, e := range errs {
			err += e.Error() + "\n"
		}
		fail(err)
	}

	return files
}

var configPath string
var excludePaths string
var formatterName string
var help bool

var originalUsage = flag.Usage

func init() {
	flag.Usage = func() {
		fmt.Println(banner)
		originalUsage()
	}
	const (
		configUsage    = "path to the configuration TOML file"
		excludeUsage   = "glob which specifies files to be excluded"
		formatterUsage = "formatter to be used for the output"
	)
	flag.StringVar(&configPath, "config", "", configUsage)
	flag.StringVar(&excludePaths, "exclude", "", excludeUsage)
	flag.StringVar(&formatterName, "formatter", "", formatterUsage)
	flag.Parse()
}
