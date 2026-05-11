package hashutil

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

// Hash hashes a blob into a hex hash that is assumed to be unique in the
// entire universe.
func Hash(bs []byte) string {
	ret := sha256.Sum256(bs)
	return hex.EncodeToString(ret[:])
}

// HashReader hashes all the bytes read from a reader.
func HashReader(r io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// HashStr hashes a string into a hash in hex.
func HashStr(s string) string {
	h := sha256.New()
	io.WriteString(h, s)
	return hex.EncodeToString(h.Sum(nil))
}

// HashFile computes the hash of the given file.
func HashFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	ret, err := HashReader(f)
	if err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}
	return ret, nil
}
