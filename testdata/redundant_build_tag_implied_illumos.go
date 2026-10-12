//go:build solaris

// MATCH:1 /The build tag "//go:build solaris" is implied by the file name "redundant_build_tag_implied_illumos.go" and can be removed/

// GOOS=illumos also satisfies the "solaris" tag.

package pkg
