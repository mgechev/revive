package lint

import "testing"

// TestName tests Name function.
func TestName(t *testing.T) {
	tests := []struct {
		name, want string
	}{
		{"foo_bar", "fooBar"},
		{"foo_bar_baz", "fooBarBaz"},
		{"Foo_bar", "FooBar"},
		{"foo_WiFi", "fooWiFi"},
		{"id", "id"},
		{"Id", "ID"},
		{"foo_id", "fooID"},
		{"fooId", "fooID"},
		{"fooUid", "fooUID"},
		{"idFoo", "idFoo"},
		{"uidFoo", "uidFoo"},
		{"midIdDle", "midIDDle"},
		{"APIProxy", "APIProxy"},
		{"ApiProxy", "APIProxy"},
		{"apiProxy", "apiProxy"},
		{"_Leading", "_Leading"},
		{"___Leading", "_Leading"},
		{"trailing_", "trailing"},
		{"trailing___", "trailing"},
		{"a_b", "aB"},
		{"a__b", "aB"},
		{"a___b", "aB"},
		{"Rpc1150", "RPC1150"},
		{"case3_1", "case3_1"},
		{"case3__1", "case3_1"},
		{"IEEE802_16bit", "IEEE802_16bit"},
		{"IEEE802_16Bit", "IEEE802_16Bit"},
	}
	for _, test := range tests {
		got := Name(test.name, nil, nil)
		if got != test.want {
			t.Errorf("Name(%q) = %q, want %q", test.name, got, test.want)
		}
	}
}
