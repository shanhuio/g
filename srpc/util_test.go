package srpc

import (
	"net/url"
	"testing"
)

func TestPathJoin(t *testing.T) {
	for _, test := range []struct {
		server, p, want string
	}{
		{server: "http://h", p: "p", want: "http://h/p"},
		{server: "http://h/", p: "p", want: "http://h/p"},
		{server: "http://h", p: "/p", want: "http://h/p"},
		{server: "http://h/", p: "/p", want: "http://h/p"},
		{server: "http://h/p", p: "q", want: "http://h/p/q"},
		{server: "http://h/p", p: "/q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "/q", want: "http://h/p/q"},
		{server: "http://h/p/", p: "q/", want: "http://h/p/q/"},
		{server: "http://h/p", p: "q/", want: "http://h/p/q/"},
	} {
		server, err := url.Parse(test.server)
		if err != nil {
			t.Fatalf("parse server %q: %s", test.server, err)
		}
		got := urlJoin(server, test.p).String()
		if got != test.want {
			t.Errorf(
				"pathJoin(%q, %q) = %q, want %q",
				test.server, test.p, got, test.want,
			)
		}
	}
}
