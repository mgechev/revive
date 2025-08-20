package fixtures

type SecretsSerializationCustom struct {
	Foo         string `json:"-"`
	Bar         string // MATCH /Struct field 'Bar' may contain secrets but is not excluded from JSON serialization (missing `json:"-"`)/
	BearerToken string
	AuthToken   string `json:"-"`
	Apikey      string
	credential  string
}
