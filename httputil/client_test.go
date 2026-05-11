package httputil

import (
	"testing"

	"io"
)

func TestClientGetCode(t *testing.T) {
	s := newHelloServer()
	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	got, err := c.GetCode("/")
	if err != nil {
		t.Fatal(err)
	}
	if got != 200 {
		t.Errorf("want 200, got %d", got)
	}

	got, err = c.GetCode("/secret")
	if err != nil {
		t.Fatal(err)
	}
	if got != 403 {
		t.Errorf("want 403, got %d", got)
	}
}

func TestClientGet(t *testing.T) {
	s := newHelloServer()
	c, err := NewClient(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Get("/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	got, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != testHelloMessage {
		t.Errorf("got %q, want %q", string(got), testHelloMessage)
	}
}
