package httputil

import (
	"testing"

	"fmt"
	"net/http"
	"net/http/httptest"
)

const testHelloMessage = "hello"

func helloHandler(w http.ResponseWriter, req *http.Request) {
	p := req.URL.Path
	if p == "/secret" {
		http.Error(w, "not authorized", 403)
		return
	}
	fmt.Fprint(w, testHelloMessage)
}

func newHelloServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(helloHandler))
}

func TestGetString(t *testing.T) {
	s := newHelloServer()
	c := s.Client()
	got, err := GetString(c, s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != testHelloMessage {
		t.Errorf("want %q, got %q", testHelloMessage, got)
	}
}

func TestGetCode(t *testing.T) {
	s := newHelloServer()
	c := s.Client()
	got, err := GetCode(c, s.URL)
	if err != nil {
		t.Fatal(err)
	}
	if got != 200 {
		t.Errorf("want 200, got %d", got)
	}

	got, err = GetCode(c, s.URL+"/secret")
	if err != nil {
		t.Fatal(err)
	}
	if got != 403 {
		t.Errorf("want 403, got %d", got)
	}
}
