package pisces

import (
	"fmt"

	"shanhu.io/g/hashutil"
)

// MaxKeyLen is the maximum length of a hashed KV.
const MaxKeyLen = 255

func keyHash(k string) string {
	return hashutil.HashStr(k)
}

func kvMapKey(key string, ordered bool) (string, error) {
	if !ordered {
		return keyHash(key), nil
	}
	if len(key) > MaxKeyLen {
		return "", fmt.Errorf("key %q too long", key)
	}
	return key, nil
}
