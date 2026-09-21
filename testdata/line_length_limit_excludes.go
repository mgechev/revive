package fixtures

//go:generate echo aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

// A reference to https://example.com/very/long/path/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa

// A perfectly ordinary but very long comment line that clearly exceeds the eighty character budget

func foo() {
	_ = "url inside code https://example.com/aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	_ = "no url here just a long enough string literal to exceed the eighty character line length yyy"
}

// MATCH:7 /line is 99 characters, out of limit 80/
// MATCH:11 /line is 102 characters, out of limit 80/
