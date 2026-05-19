package rule_test

import (
	"testing"

	"github.com/mgechev/revive/internal/rule"
)

func TestName(t *testing.T) {
	tests := []struct {
		name                     string
		allowlist                []string
		blocklist                []string
		skipInitialismNameChecks bool
		initialismsAsWords       bool
		want                     string
	}{
		{
			name: "foo_bar",
			want: "fooBar",
		},
		{
			name: "foo_bar_baz",
			want: "fooBarBaz",
		},
		{
			name: "Foo_bar",
			want: "FooBar",
		},
		{
			name: "foo_WiFi",
			want: "fooWiFi",
		},
		{
			name: "id",
			want: "id",
		},
		{
			name: "Id",
			want: "ID",
		},
		{
			name: "foo_id",
			want: "fooID",
		},
		{
			name: "fooId",
			want: "fooID",
		},
		{
			name: "fooUid",
			want: "fooUID",
		},
		{
			name: "idFoo",
			want: "idFoo",
		},
		{
			name: "uidFoo",
			want: "uidFoo",
		},
		{
			name: "midIdDle",
			want: "midIDDle",
		},
		{
			name: "APIProxy",
			want: "APIProxy",
		},
		{
			name: "ApiProxy",
			want: "APIProxy",
		},
		{
			name: "apiProxy",
			want: "apiProxy",
		},
		{
			name: "_Leading",
			want: "_Leading",
		},
		{
			name: "___Leading",
			want: "_Leading",
		},
		{
			name: "trailing_",
			want: "trailing",
		},
		{
			name: "trailing___",
			want: "trailing",
		},
		{
			name: "a_b",
			want: "aB",
		},
		{
			name: "a__b",
			want: "aB",
		},
		{
			name: "a___b",
			want: "aB",
		},
		{
			name: "Rpc1150",
			want: "RPC1150",
		},
		{
			name: "case3_1",
			want: "case3_1",
		},
		{
			name: "case3__1",
			want: "case3_1",
		},
		{
			name: "IEEE802_16bit",
			want: "IEEE802_16bit",
		},
		{
			name: "IEEE802_16Bit",
			want: "IEEE802_16Bit",
		},
		{
			name: "IDS",
			want: "IDs",
		},
		// Test skipInitialismChecks functionality
		{
			name:                     "getJson",
			skipInitialismNameChecks: true,
			want:                     "getJson",
		},
		{
			name:                     "userId",
			skipInitialismNameChecks: true,
			want:                     "userId",
		},
		{
			name:                     "myHttpClient",
			skipInitialismNameChecks: true,
			want:                     "myHttpClient",
		},
		// Test allowlist functionality
		{
			name:      "fooId",
			allowlist: []string{"ID"},
			want:      "fooId",
		},
		{
			name:      "fooApi",
			allowlist: []string{"API"},
			want:      "fooApi",
		},
		{
			name:      "fooHttp",
			allowlist: []string{"HTTP"},
			want:      "fooHttp",
		},
		// Test blocklist functionality
		{
			name:      "fooCustom",
			blocklist: []string{"CUSTOM"},
			want:      "fooCUSTOM",
		},
		{
			name:      "mySpecial",
			blocklist: []string{"SPECIAL"},
			want:      "mySPECIAL",
		},
		// Test combination of allowlist and blocklist
		{
			name:      "fooIdCustom",
			allowlist: []string{"ID"},
			blocklist: []string{"CUSTOM"},
			want:      "fooIdCUSTOM",
		},
		// Test combination of allowlist, blocklist and skipInitialismChecks
		{
			name:                     "fooIdCustomHttpJson",
			allowlist:                []string{"ID"},
			blocklist:                []string{"CUSTOM"},
			skipInitialismNameChecks: true,
			want:                     "fooIdCustomHttpJson",
		},

		{
			name:               "foo_bar",
			want:               "fooBar",
			initialismsAsWords: true,
		},
		{
			name:               "foo_bar_baz",
			want:               "fooBarBaz",
			initialismsAsWords: true,
		},
		{
			name:               "Foo_bar",
			want:               "FooBar",
			initialismsAsWords: true,
		},
		{
			name:               "foo_WiFi",
			want:               "fooWiFi",
			initialismsAsWords: true,
		},
		{
			name:               "id",
			want:               "id",
			initialismsAsWords: true,
		},
		{
			name:               "Id",
			want:               "Id",
			initialismsAsWords: true,
		},
		{
			name:               "foo_id",
			want:               "fooId",
			initialismsAsWords: true,
		},
		{
			name:               "fooId",
			want:               "fooId",
			initialismsAsWords: true,
		},
		{
			name:               "fooUid",
			want:               "fooUid",
			initialismsAsWords: true,
		},
		{
			name:               "idFoo",
			want:               "idFoo",
			initialismsAsWords: true,
		},
		{
			name:               "uidFoo",
			want:               "uidFoo",
			initialismsAsWords: true,
		},
		{
			name:               "midIdDle",
			want:               "midIdDle",
			initialismsAsWords: true,
		},
		{
			name:               "APIProxy",
			want:               "ApiProxy",
			initialismsAsWords: true,
		},
		{
			name:               "ApiProxy",
			want:               "ApiProxy",
			initialismsAsWords: true,
		},
		{
			name:               "apiProxy",
			want:               "apiProxy",
			initialismsAsWords: true,
		},
		{
			name:               "_Leading",
			want:               "_Leading",
			initialismsAsWords: true,
		},
		{
			name:               "___Leading",
			want:               "_Leading",
			initialismsAsWords: true,
		},
		{
			name:               "trailing_",
			want:               "trailing",
			initialismsAsWords: true,
		},
		{
			name:               "trailing___",
			want:               "trailing",
			initialismsAsWords: true,
		},
		{
			name:               "a_b",
			want:               "aB",
			initialismsAsWords: true,
		},
		{
			name:               "a__b",
			want:               "aB",
			initialismsAsWords: true,
		},
		{
			name:               "a___b",
			want:               "aB",
			initialismsAsWords: true,
		},
		{
			name:               "Rpc1150",
			want:               "Rpc1150",
			initialismsAsWords: true,
		},
		{
			name:               "rpc1150",
			want:               "rpc1150",
			initialismsAsWords: true,
		},
		{
			name:               "case3_1",
			want:               "case3_1",
			initialismsAsWords: true,
		},
		{
			name:               "case3__1",
			want:               "case3_1",
			initialismsAsWords: true,
		},
		{
			name:               "IEEE802_16bit",
			want:               "Ieee802_16bit",
			initialismsAsWords: true,
		},
		{
			name:               "IEEE802_16Bit",
			want:               "Ieee802_16Bit",
			initialismsAsWords: true,
		},
		{
			name:               "IDS",
			want:               "Ids",
			initialismsAsWords: true,
		},
		// Test skipInitialismChecks functionality
		{
			name:                     "getJson",
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "getJson",
		},
		{
			name:                     "userId",
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "userId",
		},
		{
			name:                     "myHttpClient",
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "myHttpClient",
		},
		// Test allowlist functionality
		{
			name:               "fooId",
			allowlist:          []string{"ID"},
			initialismsAsWords: true,
			want:               "fooId",
		},
		{
			name:               "fooApi",
			allowlist:          []string{"API"},
			initialismsAsWords: true,
			want:               "fooApi",
		},
		{
			name:               "fooHttp",
			allowlist:          []string{"HTTP"},
			initialismsAsWords: true,
			want:               "fooHttp",
		},
		// Test blocklist functionality
		{
			name:               "fooCustom",
			blocklist:          []string{"CUSTOM"},
			initialismsAsWords: true,
			want:               "fooCUSTOM",
		},
		{
			name:               "mySpecial",
			blocklist:          []string{"SPECIAL"},
			initialismsAsWords: true,
			want:               "mySPECIAL",
		},
		// Test combination of allowlist and blocklist
		{
			name:               "fooIdCustom",
			allowlist:          []string{"ID"},
			blocklist:          []string{"CUSTOM"},
			initialismsAsWords: true,
			want:               "fooIdCUSTOM",
		},
		// Test combination of allowlist, blocklist and skipInitialismChecks
		{
			name:                     "fooIdCustomHttpJson",
			allowlist:                []string{"ID"},
			blocklist:                []string{"CUSTOM"},
			skipInitialismNameChecks: true,
			initialismsAsWords:       true,
			want:                     "fooIdCustomHttpJson",
		},
	}
	for _, test := range tests {
		got := rule.Name(test.name, test.allowlist, test.blocklist, test.skipInitialismNameChecks, test.initialismsAsWords)
		if got != test.want {
			t.Errorf("name(%q, allowlist=%v, blocklist=%v, skipInitialismNameChecks=%v, initialismsAsWords=%v) = %q, want %q",
				test.name, test.allowlist, test.blocklist, test.skipInitialismNameChecks, test.initialismsAsWords, got, test.want)
		}
	}
}
