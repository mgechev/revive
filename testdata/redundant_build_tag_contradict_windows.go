//go:build !windows

// MATCH:1 /The build tag "//go:build !windows" contradicts the file name "redundant_build_tag_contradict_windows.go", so the file is never built/

package pkg
