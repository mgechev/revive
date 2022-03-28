package revivelib

// LintPattern returns either an include or exclude package
type LintPattern interface {
	IsExclude() bool
	GetPattern() string
}

// IncludePattern indicates a package that's included in the lint
type IncludePattern struct {
	Path string
}

// IsExclude package
func (p *IncludePattern) IsExclude() bool {
	return false
}

// GetPattern for this include
func (p *IncludePattern) GetPattern() string {
	return p.Path
}

// Ensure we respect the LintPattern interface
var _ LintPattern = (*IncludePattern)(nil)

// Include this path in the linter
func Include(path string) *IncludePattern {
	return &IncludePattern{
		Path: path,
	}
}

// ExcludePattern indicates a package that's included in the lint
type ExcludePattern struct {
	Path string
}

// IsExclude package
func (p *ExcludePattern) IsExclude() bool {
	return true
}

// GetPattern for this include
func (p *ExcludePattern) GetPattern() string {
	return p.Path
}

// Ensure we respect the LintPattern interface
var _ LintPattern = (*ExcludePattern)(nil)

// Exclude this path in the linter
func Exclude(path string) *ExcludePattern {
	return &ExcludePattern{
		Path: path,
	}
}
