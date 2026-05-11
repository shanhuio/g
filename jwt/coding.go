package jwt

import (
	"encoding/base64"
	"encoding/json"
)

func encodeSegmentBytes(bs []byte) string {
	return base64.RawURLEncoding.EncodeToString(bs)
}

func decodeSegmentBytes(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func encodeSegment(v interface{}) (string, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return encodeSegmentBytes(bs), nil
}

func decodeSegment(s string, v interface{}) error {
	bs, err := decodeSegmentBytes(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}
