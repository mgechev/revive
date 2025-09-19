package fixtures

import "time"

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
	BadComplex  string `validate:"gt=0,keys,eq=1|eq=2,endkeys,required"`
	BadComplex2 string `validate:"gt=0,dive,eq=1|eq=2,endkeys,required"`
	BadComplex3 string `validate:"gt=0,dive,keys,eq=1|eq=2,endkeys,endkeys,required"`
}

type TomlUser struct {
	Username string `toml:"username,omitempty"`
	Location string `toml:"location,unknown"`
}

type SpannerUserOptions struct {
	ID   int    `spanner:"user_id,mySpannerOption"`
	A    int    `spanner:"-,mySpannerOption"`       // MATCH /useless option mySpannerOption for ignored field in spanner tag/
	Name string `spanner:"full_name,unknownOption"` // MATCH /unknown option "unknownOption" in spanner tag/
}

type uselessOptions struct {
	A  int       `bson:"-,"`
	B  int       `bson:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in bson tag/
	C  int       `bson:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in bson tag/
	D  int       `datastore:"-,"`
	E  int       `datastore:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in datastore tag/
	F  int       `datastore:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in datastore tag/
	G  int       `json:"-,"`
	H  int       `json:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in json tag/
	I  int       `json:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in json tag/
	J  int       `mapstructure:"-,"`
	K  int       `mapstructure:"-,squash"`              // MATCH /useless option squash for ignored field in mapstructure tag/
	L  int       `mapstructure:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in mapstructure tag/
	M  int       `properties:"-,"`
	N  int       `properties:"-,default=15"`                           // MATCH /useless option default=15 for ignored field in properties tag/
	O  time.Time `properties:"-,layout=2006-01-02,default=2006-01-02"` // MATCH /useless options layout=2006-01-02,default=2006-01-02 for ignored field in properties tag/
	P  int       `spanner:"-,"`
	Q  int       `spanner:"-,mySpannerOption"`                 // MATCH /useless option mySpannerOption for ignored field in spanner tag/
	R  int       `spanner:"-,mySpannerOption,mySpannerOption"` // MATCH /useless options mySpannerOption,mySpannerOption for ignored field in spanner tag/
	S  int       `toml:"-,"`
	T  int       `toml:"-,omitempty"`
	U  int       `toml:"-,omitempty,omitempty"`
	V  int       `url:"-,"`
	W  int       `url:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in url tag/
	X  int       `url:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in url tag/
	Y  int       `xml:"-,"`
	Z  int       `xml:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in xml tag/
	Aa int       `xml:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in xml tag/
	Ba int       `yaml:"-,"`
	Ca int       `yaml:"-,omitempty"`           // MATCH /useless option omitempty for ignored field in yaml tag/
	Da int       `yaml:"-,omitempty,omitempty"` // MATCH /useless options omitempty,omitempty for ignored field in yaml tag/

	// MATCH:59 /unknown option "" in bson tag/
	// MATCH:62 /unknown option "" in datastore tag/
	// MATCH:68 /unknown option "" in mapstructure tag/
	// MATCH:71 /unknown or malformed option "" in properties tag/
	// MATCH:74 /unknown option "" in spanner tag/
	// MATCH:80 /unknown option "" in url tag/
	// MATCH:83 /unknown option "" in xml tag/
	// MATCH:86 /unknown option "" in yaml tag/
}
