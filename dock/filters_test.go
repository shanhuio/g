package dock

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestLabelFilters(t *testing.T) {
	for _, test := range []struct {
		label, want string
	}{
		{label: "foo=bar", want: `{"label":["foo=bar"]}`},
		{label: "", want: `{"label":[""]}`},
		{label: "key", want: `{"label":["key"]}`},
		{label: `has"quote`, want: `{"label":["has\"quote"]}`},
		{label: "unicode=日本語", want: `{"label":["unicode=日本語"]}`},
	} {
		got, err := labelFilters(test.label)
		if err != nil {
			t.Errorf("labelFilters(%q): %v", test.label, err)
			continue
		}
		if got != test.want {
			t.Errorf(
				"labelFilters(%q), got %q, want %q",
				test.label, got, test.want,
			)
		}
	}
}

func TestLabelFiltersRoundTrip(t *testing.T) {
	const label = "name=myapp"
	got, err := labelFilters(label)
	if err != nil {
		t.Fatalf("labelFilters: %v", err)
	}
	var parsed map[string][]string
	if err := json.Unmarshal([]byte(got), &parsed); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	want := map[string][]string{"label": {label}}
	if !reflect.DeepEqual(parsed, want) {
		t.Errorf("parsed = %v, want %v", parsed, want)
	}
}
