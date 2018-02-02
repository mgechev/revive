package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mgechev/revive/lint"
)

const banner = `
Welcome to:
 _ __ _____   _(_)__  _____
| '__/ _ \ \ / / \ \ / / _ \
| | |  __/\ V /| |\ V /  __/
|_|  \___| \_/ |_| \_/ \___|
`

func main() {
	config := getConfig()
	formatter := getFormatter()
	files := getFiles()

	revive := lint.New(func(file string) ([]byte, error) {
		return ioutil.ReadFile(file)
	})

	lintingRules := getLintingRules(config)

	failures, err := revive.Lint(files, lintingRules, *config)
	if err != nil {
		fail(err.Error())
	}

	formatChan := make(chan lint.Failure)
	exitChan := make(chan bool)

	var output string
	go (func() {
		output, err = formatter.Format(formatChan, config.Rules)
		if err != nil {
			fail(err.Error())
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
	if output != "" {
		fmt.Println(output)
	}

	os.Exit(exitCode)
}
