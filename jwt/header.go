package jwt

import (
	"encoding/json"

	"shanhu.io/g/errcode"
)

// Header is the JWT header.
type Header struct {
	Alg   string `json:"alg"`
	Typ   string `json:"typ"`
	KeyID string `json:"kid,omitempty"` // Key ID.
}

func (h *Header) encode() (string, error) {
	return encodeSegment(h)
}

func decodeHeader(s string) (*Header, error) {
	bs, err := decodeSegmentBytes(s)
	if err != nil {
		return nil, err
	}
	h := new(Header)
	if err := json.Unmarshal(bs, h); err != nil {
		return nil, err
	}
	return h, nil
}

func checkHeader(got, want *Header) error {
	if got.KeyID != want.KeyID {
		return errcode.InvalidArgf("kid=%q, want %q", got.KeyID, want.KeyID)
	}
	if got.Alg != want.Alg {
		return errcode.InvalidArgf("alg=%q, want %q", got.Alg, want.Alg)
	}
	if got.Typ != want.Typ {
		return errcode.InvalidArgf("typ=%q, want %q", got.Typ, want.Typ)
	}
	return nil
}
