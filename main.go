package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/mgechev/revive/formatter"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/rule"
)

var allRules = []lint.Rule{
	&rule.ArgumentsLimitRule{},
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
	&rule.CyclomaticRule{},
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
			panic("cannot find rule: " + name)
		}
		lintingRules = append(lintingRules, rule)
	}

	return lintingRules
}

func parseConfig(path string) *lint.Config {
	config := &lint.Config{}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic("cannot read the config file")
	}
	_, err = toml.Decode(string(file), config)
	if err != nil {
		panic("cannot parse the config file: " + err.Error())
	}
	return config
}

func normalizeConfig(config *lint.Config) {
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

const usage = `
Welcome to:        
 _ __ _____   _(_)__   _____ 
| '__/ _ \ \ / / \ \ / / _ \
| | |  __/\ V /| |\ V /  __/
|_|  \___| \_/ |_| \_/ \___|

Usage:
        revive [flags] <Go file or directory> ...
Flags:
        -c   string  path to the configuration TOML file.
        -e   string  glob which specifies files to be excluded.
        -f   string  formatter to be used for the output.
        -h           output this screen.
`

func main() {
	src := `
	package p

	func Test() {
		if true || bar && baz {
			return 42;
		} else {
			return 23;
		}
	}
	
	func foo_bar(a int, b int, c int, d int) {
		return a + b + c;
	}`

	revive := lint.New(func(file string) ([]byte, error) {
		return []byte(src), nil
	})

	config := parseConfig("config.toml")
	normalizeConfig(config)
	lintingRules := getLintingRules(config)

	failures, err := revive.Lint([]string{"foo.go", "bar.go", "baz.go"}, lintingRules, config.Rules)
	if err != nil {
		panic(err)
	}

	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	var output string
	go (func() {
		var formatter formatter.CLIFormatter
		output, err = formatter.Format(formatChan, config.Rules)
		if err != nil {
			panic(err)
		}
		exitChan <- true
	})()

	exitCode := 0
	for f := range failures {
		if f.Confidence < config.Confidence {
			continue
		}
		if exitCode == 0 {
			exitCode = 1
		}
		if c, ok := config.Rules[f.RuleName]; ok && c.Severity == lint.SeverityError {
			exitCode = 2
		}
		formatChan <- f
	}

	close(formatChan)
	<-exitChan
	fmt.Println(output)

	os.Exit(exitCode)
}
