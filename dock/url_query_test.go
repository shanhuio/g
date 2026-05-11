package dock

import (
	"net/url"
	"testing"
)

func TestURLQuery(t *testing.T) {
	for _, test := range []struct {
		p    string
		q    url.Values
		want string
	}{
		{p: "/v1.40/version", q: nil, want: "/v1.40/version"},
		{p: "/v1.40/version", q: url.Values{}, want: "/v1.40/version"},
		{
			p:    "/v1.40/containers/abc/kill",
			q:    url.Values{"signal": []string{"SIGINT"}},
			want: "/v1.40/containers/abc/kill?signal=SIGINT",
		},
		{
			p: "/v1.40/containers/abc/stop",
			q: url.Values{
				"t":     []string{"60"},
				"extra": []string{"a b"},
			},
			want: "/v1.40/containers/abc/stop?extra=a+b&t=60",
		},
	} {
		got := urlQuery(test.p, test.q)
		if got != test.want {
			t.Errorf(
				"urlQuery(%q, %v), got %q, want %q",
				test.p, test.q, got, test.want,
			)
		}
	}
}

func TestAPIURLQuery(t *testing.T) {
	for _, test := range []struct {
		p    string
		q    url.Values
		want string
	}{
		{p: "version", q: nil, want: "/v1.40/version"},
		{p: "/version", q: nil, want: "/v1.40/version"},
		{
			p:    "containers/create",
			q:    url.Values{"name": []string{"foo"}},
			want: "/v1.40/containers/create?name=foo",
		},
		{
			p:    "images/create",
			q:    url.Values{"fromImage": []string{"alpine"}, "tag": []string{"latest"}},
			want: "/v1.40/images/create?fromImage=alpine&tag=latest",
		},
	} {
		got := apiURLQuery(test.p, test.q)
		if got != test.want {
			t.Errorf(
				"apiURLQuery(%q, %v), got %q, want %q",
				test.p, test.q, got, test.want,
			)
		}
	}
}

func TestSingleQuery(t *testing.T) {
	for _, test := range []struct {
		k, v string
	}{
		{k: "name", v: "foo"},
		{k: "signal", v: "SIGINT"},
		{k: "t", v: "60"},
		{k: "force", v: "1"},
	} {
		q := singleQuery(test.k, test.v)
		if len(q) != 1 {
			t.Errorf(
				"singleQuery(%q, %q), got %d keys, want 1",
				test.k, test.v, len(q),
			)
		}
		if got := q.Get(test.k); got != test.v {
			t.Errorf(
				"singleQuery(%q, %q).Get(%q), got %q, want %q",
				test.k, test.v, test.k, got, test.v,
			)
		}
	}
}
