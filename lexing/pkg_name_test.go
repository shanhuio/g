package lexing

import (
	"testing"
)

func TestIsPkgName(t *testing.T) {
	testName := []string{"a", "abc", "a12", "a12b"}
	for _, name := range testName {
		if !IsPkgName(name) {
			t.Errorf("%v should be a package name", name)
		}
	}
	testName = []string{"aB", "1abc", "", "  ", "a~", "%a", "A1", "TBC", "$abc1"}
	for _, name := range testName {
		if IsPkgName(name) {
			t.Errorf("%v should not be a package name", name)
		}
	}
}
