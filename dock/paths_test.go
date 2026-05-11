package dock

import "testing"

func TestContPath(t *testing.T) {
	for _, test := range []struct {
		id, method, want string
	}{
		{id: "abc", method: "", want: "containers/abc"},
		{id: "abc", method: "start", want: "containers/abc/start"},
		{id: "abc", method: "json", want: "containers/abc/json"},
		{id: "abc", method: "archive", want: "containers/abc/archive"},
		{id: "my-name", method: "stop", want: "containers/my-name/stop"},
	} {
		got := contPath(test.id, test.method)
		if got != test.want {
			t.Errorf(
				"contPath(%q, %q), got %q, want %q",
				test.id, test.method, got, test.want,
			)
		}
	}
}

func TestExecPath(t *testing.T) {
	for _, test := range []struct {
		id, method, want string
	}{
		{id: "abc", method: "start", want: "exec/abc/start"},
		{id: "abc", method: "json", want: "exec/abc/json"},
		{id: "abc", method: "", want: "exec/abc"},
	} {
		got := execPath(test.id, test.method)
		if got != test.want {
			t.Errorf(
				"execPath(%q, %q), got %q, want %q",
				test.id, test.method, got, test.want,
			)
		}
	}
}
