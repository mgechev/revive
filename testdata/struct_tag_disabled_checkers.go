package fixtures

type DisabledCheckersTest struct {
	// These validate tags should not trigger errors since validate checker is disabled
	Field1 string `validate:"gt:0"`
	Field2 string `validate:"in:FOO,BAR,BAZ"`
	Field3 string `validate:"required_if:field1,true|max_len:255"`
	Field4 string `validate:"min_len:1|max_len:256"`
	
	// These toml tags should not trigger errors since toml checker is disabled
	Field5 string `toml:"field5,unknown_option"`
	Field6 string `toml:"field6,another_unknown"`
	
	// These json tags should still trigger errors since json checker is not disabled
	Field7 string `json:"field7,invalid_option"` // MATCH /unknown option "invalid_option" in json tag/
	Field8 string `json:"field8,omitempty"`      // This should be fine
	
	// These bson tags should still trigger errors since bson checker is not disabled  
	Field9 string `bson:"field9,invalid_option"` // MATCH /unknown option "invalid_option" in bson tag/
}