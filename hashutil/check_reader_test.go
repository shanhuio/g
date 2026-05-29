package hashutil

import (
	"testing"

	"bytes"
	"crypto/sha256"
	"io"

	"shanhu.io/std/errcode"
)

func TestCheckReader(t *testing.T) {
	msg := []byte("hello")
	h := sha256.Sum256(msg)
	msgLen := int64(len(msg))

	t.Run("good", func(t *testing.T) {
		r := bytes.NewReader(msg)
		cr := NewSHA256CheckReader(r, h[:], msgLen)
		if _, err := io.Copy(io.Discard, cr); err != nil {
			t.Errorf("want no error, got %q", err)
		}
	})

	t.Run("noLenCheck", func(t *testing.T) {
		r := bytes.NewReader(msg)
		cr := NewSHA256CheckReader(r, h[:], -1)
		if _, err := io.Copy(io.Discard, cr); err != nil {
			t.Errorf("want no error, got %q", err)
		}
	})

	t.Run("wrongHash", func(t *testing.T) {
		h2 := make([]byte, len(h))
		copy(h2, h[:])
		h2[0] = ^h2[0] // change the hash
		r := bytes.NewReader(msg)
		cr := NewSHA256CheckReader(r, h2[:], msgLen)
		if _, err := io.Copy(io.Discard, cr); !errcode.IsInvalidArg(err) {
			t.Errorf("want invalid arg, got %q", err)
		}
	})

	t.Run("wrongLen", func(t *testing.T) {
		r := bytes.NewReader(msg)
		cr := NewSHA256CheckReader(r, h[:], msgLen+1)
		if _, err := io.Copy(io.Discard, cr); !errcode.IsInvalidArg(err) {
			t.Errorf("want invalid arg, got %q", err)
		}
	})
}
