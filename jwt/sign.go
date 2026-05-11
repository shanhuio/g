package jwt

import (
	"bytes"
	"context"
	"io"

	"shanhu.io/g/errcode"
)

// Signer signs the token, returns the signature and the header.
type Signer interface {
	Header(ctx context.Context) (*Header, error)
	Sign(ctx context.Context, h *Header, data []byte) ([]byte, error)
}

// EncodeAndSign signs and encodes a claim set and signs it.
func EncodeAndSign(ctx context.Context, c *ClaimSet, s Signer) (string, error) {
	h, err := s.Header(ctx)
	if err != nil {
		return "", errcode.Annotate(err, "get header")
	}
	hb, err := h.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode header")
	}

	cb, err := c.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode claims")
	}
	buf := new(bytes.Buffer)
	io.WriteString(buf, hb)
	io.WriteString(buf, ".")
	io.WriteString(buf, cb)
	sig, err := s.Sign(ctx, h, buf.Bytes())
	if err != nil {
		return "", errcode.Annotate(err, "signing token")
	}
	io.WriteString(buf, ".")
	io.WriteString(buf, encodeSegmentBytes(sig))
	return buf.String(), nil
}
