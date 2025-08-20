package fixtures

type SecretsSerializationDefault struct {
	Foo         string
	Bar         string
	BearerToken string // MATCH /Struct field 'BearerToken' may contain secrets but is not excluded from JSON serialization (missing `json:"-"`)/
	AuthToken   string `json:"-"`
	Apikey      string //MATCH /Struct field 'Apikey' may contain secrets but is not excluded from JSON serialization (missing `json:"-"`)/
	credential  string
}
