package timeutil

import (
	"encoding/base64"
	"io"
	"time"

	"shanhu.io/std/errcode"
)

// Challenge is a timestamp with a crypto random nonce. A server can provide
// an HMAC signed challenge for clients via RPC, and use it as a structure
// to restrict the client into a time window that is defined by the server.
// It also provides a reliable clock source and a nonce source. The nonce
// can be used as a one-time token to defend replay attacks.
type Challenge struct {
	N string // Nonce.
	T *Timestamp
}

// NewChallenge creates a challenge with the given timestamp t.
// rand is used to generate the nonce.
func NewChallenge(t time.Time, rand io.Reader) (*Challenge, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, errcode.Annotate(err, "read nonce")
	}

	nonceStr := base64.RawStdEncoding.EncodeToString(nonce)
	return &Challenge{
		N: nonceStr,
		T: NewTimestamp(t),
	}, nil
}
