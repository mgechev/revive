// Package analyzer provides revive as a golang.org/x/tools/go/analysis analyzer.
package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	goversion "github.com/hashicorp/go-version"
	"golang.org/x/tools/go/analysis"

	"github.com/mgechev/revive/config"
	"github.com/mgechev/revive/lint"
	"github.com/mgechev/revive/revivelib"
)

// New returns an analyzer applying the rules enabled by conf and the given extra rules.
func New(conf *lint.Config, extraRules ...revivelib.ExtraRule) (*analysis.Analyzer, error) {
	extraRuleInstances := make([]lint.Rule, len(extraRules))
	for i, extraRule := range extraRules {
		extraRuleInstances[i] = extraRule.Rule

		ruleName := extraRule.Rule.Name()
		if _, ok := conf.Rules[ruleName]; !ok {
			conf.Rules[ruleName] = extraRule.DefaultConfig
		}
	}

	rules, err := config.GetLintingRules(conf, extraRuleInstances)
	if err != nil {
		return nil, fmt.Errorf("creating analyzer - getting lint rules: %w", err)
	}

	return &analysis.Analyzer{
		Name: "revive",
		Doc:  "fast, configurable, extensible, flexible, and beautiful linter for Go",
		Run: func(pass *analysis.Pass) (any, error) {
			return nil, run(pass, *conf, rules)
		},
	}, nil
}

func run(pass *analysis.Pass, conf lint.Config, rules []lint.Rule) error {
	pkg := lint.NewPackage(pass.Fset, goVersion(pass, conf), pass.Pkg, pass.TypesInfo)

	tokenFiles := map[string]*token.File{}
	for _, astFile := range pass.Files {
		if !conf.IgnoreGeneratedHeader && ast.IsGenerated(astFile) {
			continue
		}

		tokenFile := pass.Fset.File(astFile.Pos())
		if tokenFile == nil {
			continue
		}

		filename := tokenFile.Name()
		content, err := pass.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("reading file %q: %w", filename, err)
		}

		pkg.AddFile(lint.NewFileFromAST(filename, content, astFile, pkg))
		tokenFiles[filename] = tokenFile
	}

	failures := make(chan lint.Failure)
	lintErr := make(chan error, 1)
	go func() {
		lintErr <- pkg.Lint(rules, conf, failures)
		close(failures)
	}()

	for failure := range failures {
		pass.Report(toDiagnostic(failure, tokenFiles))
	}

	return <-lintErr
}

func goVersion(pass *analysis.Pass, conf lint.Config) *goversion.Version {
	if conf.GoVersion != nil {
		return conf.GoVersion
	}

	version := strings.TrimPrefix(pass.Module.GoVersion, "go")
	if version == "" {
		return nil
	}

	return goversion.Must(goversion.NewVersion(version))
}

func toDiagnostic(failure lint.Failure, tokenFiles map[string]*token.File) analysis.Diagnostic {
	diagnostic := analysis.Diagnostic{
		Category: failure.RuleName,
		Message:  failure.Failure,
	}

	if failure.Node != nil {
		diagnostic.Pos = failure.Node.Pos()
		diagnostic.End = failure.Node.End()
	} else {
		diagnostic.Pos = toPos(failure.Position.Start, tokenFiles)
		diagnostic.End = toPos(failure.Position.End, tokenFiles)
	}

	if fix := toSuggestedFix(failure, tokenFiles); fix != nil {
		diagnostic.SuggestedFixes = []analysis.SuggestedFix{*fix}
	}

	return diagnostic
}

func toPos(position token.Position, tokenFiles map[string]*token.File) token.Pos {
	tokenFile, ok := tokenFiles[position.Filename]
	if !ok || position.Line < 1 || position.Line > tokenFile.LineCount() {
		return token.NoPos
	}

	pos := tokenFile.LineStart(position.Line)
	if position.Column > 1 {
		pos += token.Pos(position.Column - 1)
	}

	return min(pos, fileEnd(tokenFile))
}

func toSuggestedFix(failure lint.Failure, tokenFiles map[string]*token.File) *analysis.SuggestedFix {
	if failure.ReplacementLine == "" {
		return nil
	}

	position := failure.Position.Start
	tokenFile, ok := tokenFiles[position.Filename]
	if !ok || position.Line < 1 || position.Line > tokenFile.LineCount() {
		return nil
	}

	lineStart := tokenFile.LineStart(position.Line)
	lineEnd := fileEnd(tokenFile)
	if position.Line < tokenFile.LineCount() {
		lineEnd = tokenFile.LineStart(position.Line+1) - 1
	}

	return &analysis.SuggestedFix{
		Message: "Replace the line by: " + failure.ReplacementLine,
		TextEdits: []analysis.TextEdit{{
			Pos:     lineStart,
			End:     lineEnd,
			NewText: []byte(failure.ReplacementLine),
		}},
	}
}

func fileEnd(tokenFile *token.File) token.Pos {
	return token.Pos(tokenFile.Base() + tokenFile.Size())
}
