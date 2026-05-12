package fixtures

import "encoding/json"

// Good: MarshalJSON with value receiver.
type GoodMarshalRecv struct{}

func (GoodMarshalRecv) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) }
func (GoodMarshalRecv) MarshalText() (text []byte, err error)  { return json.Marshal(nil) }
func (*GoodMarshalRecv) UnmarshalJSON([]byte) error   { return nil }
func (GoodMarshalRecv) MarshalYAML() (any, error)     { var v any; return v, nil }
func (*GoodMarshalRecv) UnmarshalYAML(func(any) error) error { return nil }

// Bad: MarshalJSON with pointer receiver.
type BadMarshalRecv struct{}

func (*BadMarshalRecv) MarshalJSON() ([]byte, error)  { return json.Marshal(nil) } // MATCH /MarshalJSON method should use a value receiver, not a pointer receiver/
func (BadMarshalRecv) UnmarshalText(text []byte) error     { return nil }               // MATCH /UnmarshalText method should use a pointer receiver, not a value receiver/
func (BadMarshalRecv) UnmarshalJSON([]byte) error     { return nil }               // MATCH /UnmarshalJSON method should use a pointer receiver, not a value receiver/
func (*BadMarshalRecv) MarshalYAML() (any, error)     { var v any; return v, nil }          // MATCH /MarshalYAML method should use a value receiver, not a pointer receiver/
func (BadMarshalRecv) UnmarshalYAML(func(any) error) error { return nil }          // MATCH /UnmarshalYAML method should use a pointer receiver, not a value receiver/
