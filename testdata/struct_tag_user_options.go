package fixtures

type RangeAllocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Range             string `json:"range,outline"`
	Data              []byte `json:"data,flow"` // MATCH /unknown option "flow" in json tag/
}

type RangeAllocation struct {
	metav1.TypeMeta   `bson:",minsize,gnu"`
	metav1.ObjectMeta `bson:"metadata,omitempty"`
	Range             string `bson:"range,flow"` // MATCH /unknown option "flow" in bson tag/
	Data              []byte `bson:"data,inline"`
}

type RequestQueryOptions struct {
	Properties       []string `url:"properties,commmma,omitempty"` // MATCH /unknown option "commmma" in url tag/
	CustomProperties []string `url:"-"`
	Archived         bool     `url:"archived,myURLOption"`
}

type Fields struct {
	Field      string `datastore:",noindex,flatten,omitempty,myDatastoreOption"`
	OtherField string `datastore:",unknownOption"` // MATCH /unknown option "unknownOption" in datastore tag/
}

type MapStruct struct {
	Field1     string `mapstructure:",squash,reminder,omitempty,myMapstructureOption"`
	OtherField string `mapstructure:",unknownOption"` // MATCH /unknown option "unknownOption" in mapstructure tag/
}

type ValidateUser struct {
	Username    string `validate:"required,min=3,max=32"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required,min=8,max=32"`
	Biography   string `validate:"min=0,max=1000"`
	DisplayName string `validate:"displayName,min=3,max=32"`
	Complex     string `validate:"gt=0,dive,keys,eq=1|eq=2,endkeys,required"`
	BadComplex  string `validate:"gt=0,keys,eq=1|eq=2,endkeys,required"`              // MATCH /option "keys" must follow a "dive" option in validate tag/
	BadComplex2 string `validate:"gt=0,dive,eq=1|eq=2,endkeys,required"`              // MATCH /option "endkeys" without a previous "keys" option in validate tag/
	BadComplex3 string `validate:"gt=0,dive,keys,eq=1|eq=2,endkeys,endkeys,required"` // MATCH /option "endkeys" without a previous "keys" option in validate tag/
}

type TomlUser struct {
	Username string `toml:"username,omitempty"`
	Location string `toml:"location,unknown"`
}

type SpannerUserOptions struct {
	ID   int    `spanner:"user_id,mySpannerOption"`
	A    int    `spanner:"-,mySpannerOption"`       // MATCH /useless option(s) mySpannerOption for ignored field in spanner tag/
	Name string `spanner:"full_name,unknownOption"` // MATCH /unknown option "unknownOption" in spanner tag/
}
