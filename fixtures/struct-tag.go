package fixtures

type decodeAndValidateRequest struct {
	// BEAWRE : the flag of URLParam should match the const string URLParam
	URLParam         string          `json:"-" path:"url_param" validate:"numeric"`
	Text             string          `json:"text" validate:"max=10"`
	DefaultInt       int             `json:"defaultInt" default:"10.0"` // MATCH /field's type and default value's type mismatch/
	DefaultInt2      int             `json:"defaultInt" default:"10"`
	DefaultString    string          `json:"defaultString" default:"foo"`
	DefaultBool      bool            `json:"defaultBool" default:"trues"` // MATCH /field's type and default value's type mismatch/
	DefaultBool2     bool            `json:"defaultBool" default:"true"`
	DefaultBool3     bool            `json:"defaultBool" default:"false"`
	DefaultFloat     float64         `json:"defaultFloat" default:"f10.0"` // MATCH /field's type and default value's type mismatch/
	DefaultFloat2    float64         `json:"defaultFloat" default:"10.0"`
	MandatoryStruct  mandatoryStruct `json:"mandatoryStruct" required:"trues"` // MATCH /required should be 'true' or 'false'/
	MandatoryStruct2 mandatoryStruct `json:"mandatoryStruct" required:"true"`
	MandatoryStruct4 mandatoryStruct `json:"mandatoryStruct" required:"false"`
	OptionalStruct   *optionalStruct `json:"optionalStruct,omitempty"`
	OptionalQuery    string          `json:"-" querystring:"queryfoo"`
}

type RangeAllocation struct {
	metav1.TypeMeta   `json:",inline"` // MATCH /unknown option 'inline' in JSON tag/
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Range             string `json:"range,flow"`  // MATCH /unknown option 'flow' in JSON tag/
	Data              []byte `json:"data,inline"` // MATCH /unknown option 'inline' in JSON tag/
}

type RangeAllocation struct {
	metav1.TypeMeta   `bson:",minsize"`
	metav1.ObjectMeta `bson:"metadata,omitempty"`
	Range             string `bson:"range,flow"` // MATCH /unknown option 'flow' in BSON tag/
	Data              []byte `bson:"data,inline"`
}

type TestContextSpecificTags2 struct {
	A       int       `asn1:"explicit,tag:1"`
	B       int       `asn1:"tag:2"`
	S       string    `asn1:"tag:0,utf8"`
	Ints    []int     `asn1:"set"`
	Version int       `asn1:"optional,explicit,default:0,tag:0"` // MATCH /duplicated tag number 0/
	Time    time.Time `asn1:"explicit,tag:4,other"`              // MATCH /unknown option 'other' in ASN1 tag/
}

type VirtualMachineRelocateSpecDiskLocator struct {
	DynamicData

	DiskId          int32                           `xml:"diskId,attr,cdata"`
	Datastore       ManagedObjectReference          `xml:"datastore,chardata,innerxml"`
	DiskMoveType    string                          `xml:"diskMoveType,omitempty,comment"`
	DiskBackingInfo BaseVirtualDeviceBackingInfo    `xml:"diskBackingInfo,omitempty,any"`
	Profile         []BaseVirtualMachineProfileSpec `xml:"profile,omitempty,other"` // MATCH /unknown option 'other' in XML tag/
}
