package fixtures

type SecretsSerializationCustom struct {
	Ssn         string `json:"-"`
	Email       string // MATCH /Struct field 'Email' may contain secrets but is not excluded from JSON serialization (missing `json:"-"`)/
	BearerToken string
	AuthToken   string `json:"-"`
	Apikey      string
	credential  string
}
