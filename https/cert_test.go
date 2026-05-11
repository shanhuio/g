package https

import (
	"testing"

	"shanhu.io/g/errcode"
)

func TestNewCACert(t *testing.T) {
	cert, err := NewCACert("test.shanhu.io")
	if err != nil {
		t.Fatalf("NewCACert() got error: %s", err)
	}

	if _, err := cert.X509KeyPair(); err != nil {
		t.Fatalf("convert to tls cert got error: %s", err)
	}
}

func TestMakeRSACertWithNoHost(t *testing.T) {
	_, err := MakeRSACert(&CertConfig{}, 0)
	if !errcode.IsInvalidArg(err) {
		t.Errorf("expect invalid arg error without hosts, got %s", err)
	}
}
