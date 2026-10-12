//go:build linux

// +build linux

// MATCH:3 /The build tag "// +build" is redundant since Go 1.17 and can be removed/

// The two tags need not share a comment group: the `// +build` one is still the only actionable
// finding, and the `//go:build` one is not reported as implied by the file name.

package pkg
