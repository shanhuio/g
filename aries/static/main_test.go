package static

import (
	"testing"

	"shanhu.io/g/aries/ariestest"
	"shanhu.io/g/httputil"
)

func TestMain(t *testing.T) {
	service := makeService("testdata")
	s, err := ariestest.HTTPSServer("shanhu.io", service)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := httputil.NewClientMust("https://shanhu.io")
	c.Transport = s.Transport

	str, err := c.GetString("/")
	if err != nil {
		t.Fatal(err)
	}
	const want = "hello\n"
	if str != want {
		t.Errorf("get /, want %q, got %q", want, str)
	}
}
