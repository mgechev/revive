package fixtures

type decodeAndValidateRequest struct {
	// BEWARE : the flag of URLParam should match the const string URLParam
	URLParam    string `json:"-" path:"url_param" validate:"numeric"`
	Text        string `json:"text" validate:"max=10"`
	DefaultInt  int    `json:"defaultInt" default:"10.0"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultInt2 int    `json:"defaultInt2" default:"10"`
	// MATCH:10 /unknown option "inline" in json tag/
	DefaultInt3      int             `json:"defaultInt2,inline" default:"11"` // MATCH /duplicated tag name "defaultInt2" in json tag/
	DefaultString    string          `json:"defaultString" default:"foo"`
	DefaultBool      bool            `json:"defaultBool" default:"trues"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultBool2     bool            `json:"defaultBool2" default:"true"`
	DefaultBool3     bool            `json:"defaultBool3" default:"false"`
	DefaultFloat     float64         `json:"defaultFloat" default:"f10.0"` // MATCH /type mismatch between field type and default value type in default tag/
	DefaultFloat2    float64         `json:"defaultFloat2" default:"10.0"`
	MandatoryStruct  mandatoryStruct `json:"mandatoryStruct" required:"trues"` // MATCH /required should be "true" or "false" in required tag/
	MandatoryStruct2 mandatoryStruct `json:"mandatoryStruct2" required:"true"`
	MandatoryStruct4 mandatoryStruct `json:"mandatoryStruct4" required:"false"`
	OptionalStruct   *optionalStruct `json:"optionalStruct,omitempty"`
	OptionalQuery    string          `json:"-" querystring:"queryfoo"`
	optionalQuery    string          `json:"-" querystring:"queryfoo"` // MATCH /tag on not-exported field optionalQuery/
	// No-reg test for bug https://github.com/mgechev/revive/issues/208
	Tiret       string `json:"-,"`                   // MATCH /useless empty option for ignored field (remove the comma after -) in json tag/
	BadTiret    string `json:"other,"`               // MATCH /option can not be empty in json tag/
	ForOmitzero string `json:"forOmitZero,omitzero"` // Go 1.24 introduces omitzero
}

type RangeAllocation struct {
	metav1.TypeMeta   `json:",inline"` // MATCH /unknown option "inline" in json tag/
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Range             string `json:"range,flow"`  // MATCH /unknown option "flow" in json tag/
	Data              []byte `json:"data,inline"` // MATCH /unknown option "inline" in json tag/
}
