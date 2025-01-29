// Test for docs in const blocks

// Package foo ...
package foo

const (
	InlineComment = "ShouldBeOK" // InlineComment is not a valid documentation
	// MATCH:7 /exported const InlineComment should have comment (or a comment on this block) or be unexported/

	// MATCH /comment on exported const InlineWhatever should be of the form "InlineWhatever ..."/
	InlineWhatever = "blah"

	Whatsit = "missing_comment"

	// We should only warn once per block for missing comments,
	// thus do not warn on:13 /exported const Whatsit should have comment (or a comment on this block) or be unexported/
	// but always complain about malformed comments.

	WhosYourDaddy = "another_missing_one"

	// Something
	// MATCH /comment on exported const WhatDoesHeDo should be of the form "WhatDoesHeDo ..."/
	WhatDoesHeDo = "it's not a tumor!"
)

// These shouldn't need doc comments.
const (
	Alpha = "a"
	Beta  = "b"
	Gamma = "g"
)

// The comment on the previous const block shouldn't flow through to here.

const UndocAgain = 6

// MATCH:35 /exported const UndocAgain should have comment or be unexported/

const (
	SomeUndocumented = 7
	// MATCH:40 /exported const SomeUndocumented should have comment (or a comment on this block) or be unexported/
)
