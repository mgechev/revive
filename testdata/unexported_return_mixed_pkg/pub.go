package pub

type impl[T any] struct{ val T }

func New[T any](v T) *impl[T] { return &impl[T]{val: v} } // MATCH /exported func New returns unexported type *pub.impl[T], which can be annoying to use/
