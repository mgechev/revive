// Test of detached package comment.

/*
Package foo is pretty sweet.
*/

package foo

// MATCH:6 /package comment is detached; there should be no blank lines between it and the package statement/
