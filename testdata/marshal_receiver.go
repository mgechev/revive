package fixtures

import "encoding/json"

// Good: MarshalJSON with value receiver.
type A struct{}

func (A) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) }
func (*A) UnmarshalJSON([]byte) error   { return nil }
func (A) MarshalYAML() (any, error)     { return nil, nil }
func (*A) UnmarshalYAML(func(any) error) error { return nil }

// Bad: MarshalJSON with pointer receiver.
type B struct{}

func (*B) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) } // MATCH /MarshalJSON method should use a value receiver, not a pointer receiver/
func (B) UnmarshalJSON([]byte) error     { return nil }               // MATCH /UnmarshalJSON method should use a pointer receiver, not a value receiver/
func (*B) MarshalYAML() (any, error)     { return nil, nil }          // MATCH /MarshalYAML method should use a value receiver, not a pointer receiver/
func (B) UnmarshalYAML(func(any) error) error { return nil }          // MATCH /UnmarshalYAML method should use a pointer receiver, not a value receiver/
