package states

import (
	"strings"

	"shanhu.io/std/errcode"
)

// GetToken gets a token string from the given key. The fetched value is
// treated as a string and whitespaces are trimmed.
func GetToken(ctx C, s States, key string) (string, error) {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// GetTokenDefault gets a token string from the given keky. The fetched value
// is treated as a string and whitespaces are trimmed. If the key does not
// exist, v is returned.
func GetTokenDefault(ctx C, s States, key, v string) (string, error) {
	tok, err := GetToken(ctx, s, key)
	if errcode.IsNotFound(err) {
		return v, nil
	}
	return tok, err
}
