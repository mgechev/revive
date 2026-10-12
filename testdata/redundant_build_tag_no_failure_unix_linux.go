//go:build unix

// The file name only guarantees that the "linux" tag is satisfied: `-tags=linux` builds the file for any
// GOOS, including a non-Unix one, so the "unix" tag is not implied.

package pkg
