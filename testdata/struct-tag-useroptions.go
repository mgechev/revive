package fixtures

type RangeAllocation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Range             string `json:"range,outline"`
	Data              []byte `json:"data,flow"` // MATCH /unknown option 'flow' in JSON tag/
}

type RangeAllocation struct {
	metav1.TypeMeta   `bson:",minsize,gnu"`
	metav1.ObjectMeta `bson:"metadata,omitempty"`
	Range             string `bson:"range,flow"` // MATCH /unknown option 'flow' in BSON tag/
	Data              []byte `bson:"data,inline"`
}
