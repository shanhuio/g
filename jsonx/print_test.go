package jsonx

import (
	"testing"

	"bytes"
	"encoding/json"
)

func TestMarshal_loopback(t *testing.T) {
	for _, obj := range []interface{}{
		"something",
		1.234,
		1234,
		nil,
		struct{ A, B string }{A: "a", B: "b"},
		map[string]string{
			"a.com": "a:8888",
			"b.com": "b:7777",
		},
		[]int{1, 2, 3},
	} {
		want, err := json.Marshal(obj)
		if err != nil {
			t.Fatalf("marshal %v: %v", obj, err)
		}

		x, err := Marshal(obj)
		if err != nil {
			t.Errorf("format %v: %v", obj, err)
			continue
		}

		var box interface{}
		if err := Unmarshal(x, &box); err != nil {
			t.Errorf("unmarshal %q: %v", x, err)
			continue
		}

		got, err := json.Marshal(box)
		if err != nil {
			t.Fatalf("marshal jsonx-gen %v: %v", obj, err)
		}

		if !bytes.Equal(want, got) {
			t.Errorf(
				"format test failed %v: got %q, want %q",
				obj, got, want,
			)
		}
	}
}
