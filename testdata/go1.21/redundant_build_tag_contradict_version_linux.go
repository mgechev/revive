//go:build linux && !go1.21

// MATCH:1 /The build tag "//go:build linux && !go1.21" contradicts the go directive in go.mod, so the file is never built/

package fixtures
