package pathutil

import (
	"testing"

	"reflect"
)

func TestSplit(t *testing.T) {
	for _, test := range []struct {
		p     string
		parts []string
	}{
		{"a", []string{"a"}},
		{"a/b", []string{"a", "b"}},
		{"shanhu.io/smlvm", []string{"shanhu.io", "smlvm"}},
	} {
		parts, err := Split(test.p)
		if err != nil {
			t.Errorf("Split(%q) got error: %s", test.p, err)
			continue
		}

		if !reflect.DeepEqual(parts, test.parts) {
			t.Errorf(
				"Split(%q) want %v got %v",
				test.p, test.parts, parts,
			)
		}
	}
}

func TestSplitInvalidPath(t *testing.T) {
	for _, p := range []string{
		"",
		"/x",
		"x//y",
		"/",
		"a/b/c/",
	} {
		parts, err := Split(p)
		if err == nil {
			t.Errorf(
				"split path %q got parts %v, want error",
				p, parts,
			)
		}
	}
}
