package rsautil

import (
	"testing"

	"os"
	"path/filepath"
	"reflect"
)

func TestReadKey(t *testing.T) {
	tmp, err := os.MkdirTemp("", "rsautil")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Recreate the test.pem key file to make sure it has the right permission
	// bits.
	testPem, err := os.ReadFile("testdata/test.pem")
	if err != nil {
		t.Fatal("read test key content: ", err)
	}

	privateKeyFile := filepath.Join(tmp, "test.pem")
	if err := os.WriteFile(privateKeyFile, testPem, 0600); err != nil {
		t.Fatal("create test key file: ", err)
	}

	privateKey, err := ReadPrivateKey(privateKeyFile)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ReadPublicKey("testdata/test.pub")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&privateKey.PublicKey, publicKey) {
		t.Error("public/private key pair not matching")
	}
}
