package revivelib

// LintPattern indicates a pattern to be included/excluded when linting
type LintPattern struct {
	isExclude bool
	pattern   string
}

// IsExclude - should this pattern be included or excluded when linting
func (p *LintPattern) IsExclude() bool {
	return p.isExclude
}

// GetPattern - returns the actual pattern
func (p *LintPattern) GetPattern() string {
	return p.pattern
}

// Include this pattern when linting
func Include(pattern string) *LintPattern {
	return &LintPattern{
		isExclude: false,
		pattern:   pattern,
	}
}

// Exclude this pattern when linting
func Exclude(pattern string) *LintPattern {
	return &LintPattern{
		isExclude: true,
		pattern:   pattern,
	}
}
