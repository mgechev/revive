package lint

// Arguments is type used for the arguments of a rule.
type Arguments = []interface{}

type FileFilters = []*FileFilter

// RuleConfig is type used for the rule configuration.
type RuleConfig struct {
	Arguments Arguments
	Severity  Severity
	Disabled  bool
	// Exclude - rule-level file excludes, TOML related (strings)
	Exclude []string
	// excludeFilters - regex-based file filters, initialized from Exclude
	excludeFilters []*FileFilter
}

// Initialize - should be called after reading from TOML file
func (rc *RuleConfig) Initialize() error {
	for _, f := range rc.Exclude {
		ff, err := ParseFileFilter(f)
		if err != nil {
			return err
		}
		rc.excludeFilters = append(rc.excludeFilters, ff)
	}
	return nil
}

// RulesConfig defines the config for all rules.
type RulesConfig = map[string]RuleConfig

// Match - checks if given [File] `f` should be covered with configured rule (not excluded)
func (rcfg *RuleConfig) Match(f *File) bool {
	for _, exclude := range rcfg.excludeFilters {
		if exclude.Match(f) {
			return false
		}
	}
	return true
}

// DirectiveConfig is type used for the linter directive configuration.
type DirectiveConfig struct {
	Severity Severity
}

// DirectivesConfig defines the config for all directives.
type DirectivesConfig = map[string]DirectiveConfig

// Config defines the config of the linter.
type Config struct {
	IgnoreGeneratedHeader bool `toml:"ignoreGeneratedHeader"`
	Confidence            float64
	Severity              Severity
	EnableAllRules        bool             `toml:"enableAllRules"`
	Rules                 RulesConfig      `toml:"rule"`
	ErrorCode             int              `toml:"errorCode"`
	WarningCode           int              `toml:"warningCode"`
	Directives            DirectivesConfig `toml:"directive"`
	Exclude               []string         `toml:"exclude"`
}
