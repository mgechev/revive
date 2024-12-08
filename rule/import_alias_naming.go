package rule

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/mgechev/revive/lint"
)

// ImportAliasNamingRule lints import alias naming.
type ImportAliasNamingRule struct {
	allowRegexp *regexp.Regexp
	denyRegexp  *regexp.Regexp

	configureOnce sync.Once
}

const defaultImportAliasNamingAllowRule = "^[a-z][a-z0-9]{0,}$"

var defaultImportAliasNamingAllowRegexp = regexp.MustCompile(defaultImportAliasNamingAllowRule)

func (r *ImportAliasNamingRule) configure(arguments lint.Arguments) error {
	if len(arguments) == 0 {
		r.allowRegexp = defaultImportAliasNamingAllowRegexp
		return nil
	}

	switch namingRule := arguments[0].(type) {
	case string:
		err := r.setAllowRule(namingRule)
		if err != nil {
			return err
		}
	case map[string]any: // expecting map[string]string
		for k, v := range namingRule {
			switch k {
			case "allowRegex":
				err := r.setAllowRule(v)
				if err != nil {
					return err
				}
			case "denyRegex":
				err := r.setDenyRule(v)
				if err != nil {
					return err
				}

			default:
				return fmt.Errorf("Invalid map key for 'import-alias-naming' rule. Expecting 'allowRegex' or 'denyRegex', got %v", k)
			}
		}
	default:
		return fmt.Errorf("invalid argument '%v' for 'import-alias-naming' rule. Expecting string or map[string]string, got %T", arguments[0], arguments[0])
	}

	if r.allowRegexp == nil && r.denyRegexp == nil {
		r.allowRegexp = defaultImportAliasNamingAllowRegexp
	}
	return nil
}

// Apply applies the rule to given file.
func (r *ImportAliasNamingRule) Apply(file *lint.File, arguments lint.Arguments) []lint.Failure {
	var configureErr error
	r.configureOnce.Do(func() { configureErr = r.configure(arguments) })

	if configureErr != nil {
		return []lint.Failure{lint.NewInternalFailure(configureErr.Error())}
	}

	var failures []lint.Failure

	for _, is := range file.AST.Imports {
		path := is.Path
		if path == nil {
			continue
		}

		alias := is.Name
		if alias == nil || alias.Name == "_" || alias.Name == "." { // "_" and "." are special types of import aliases and should be processed by another linter rule
			continue
		}

		if r.allowRegexp != nil && !r.allowRegexp.MatchString(alias.Name) {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("import name (%s) must match the regular expression: %s", alias.Name, r.allowRegexp.String()),
				Node:       alias,
				Category:   "imports",
			})
		}

		if r.denyRegexp != nil && r.denyRegexp.MatchString(alias.Name) {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Failure:    fmt.Sprintf("import name (%s) must NOT match the regular expression: %s", alias.Name, r.denyRegexp.String()),
				Node:       alias,
				Category:   "imports",
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*ImportAliasNamingRule) Name() string {
	return "import-alias-naming"
}

func (r *ImportAliasNamingRule) setAllowRule(value any) error {
	namingRule, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid argument '%v' for import-alias-naming allowRegexp rule. Expecting string, got %T", value, value)
	}

	namingRuleRegexp, err := regexp.Compile(namingRule)
	if err != nil {
		return fmt.Errorf("invalid argument to the import-alias-naming allowRegexp rule. Expecting %q to be a valid regular expression, got: %v", namingRule, err)
	}
	r.allowRegexp = namingRuleRegexp
	return nil
}

func (r *ImportAliasNamingRule) setDenyRule(value any) error {
	namingRule, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid argument '%v' for import-alias-naming denyRegexp rule. Expecting string, got %T", value, value)
	}

	namingRuleRegexp, err := regexp.Compile(namingRule)
	if err != nil {
		return fmt.Errorf("invalid argument to the import-alias-naming denyRegexp rule. Expecting %q to be a valid regular expression, got: %v", namingRule, err)
	}
	r.denyRegexp = namingRuleRegexp
	return nil
}
