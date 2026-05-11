package jsonx

import (
	"testing"
)

func TestToJSON(t *testing.T) {
	for _, test := range []struct {
		in, out string
	}{
		{`1234`, `1234`},
		{`true`, `true`},
		{`false`, `false`},
		{`null`, `null`},
		{`{value:42}`, `{"value":42}`},
		{`{value:null}`, `{"value":null}`},
		{`{bool:false}`, `{"bool":false}`},
		{`{a:42,b:true}`, `{"a":42,"b":true}`},
		{`{a:42,}`, `{"a":42}`},
		{`{a:-42}`, `{"a":-42}`},
		{`{a:+42}`, `{"a":42}`},
		{`"string"`, `"string"`},
		{"`string\n`", `"string\n"`},
		{"{a:42,\n}", `{"a":42}`},
		{"{a:{a:{a:42}}}", `{"a":{"a":{"a":42}}}`},
		{"{\na:42,\n}", `{"a":42}`},
		{"42 // comment", "42"},
		{`{}`, `{}`},
		{"{a:/*a*/42}", `{"a":42}`},
		{`{"a":"a","b":"b"}`, `{"a":"a","b":"b"}`},
		{"a.b.c.d", `["a","b","c","d"]`},
	} {
		bs, errs := ToJSON([]byte(test.in))
		if errs != nil {
			t.Errorf("convert %q, got error %q", test.in, errs[0])
			continue
		}
		if got := string(bs); got != test.out {
			t.Errorf("convert %q, got %q, want %q", test.in, got, test.out)
		}
	}
}
