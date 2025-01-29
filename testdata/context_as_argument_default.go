package foo

import (
	"context"
	"testing"
)

func x(_ AllowedBeforePtrStruct, ctx context.Context) { // MATCH /context.Context should be the first parameter of a function/
}
