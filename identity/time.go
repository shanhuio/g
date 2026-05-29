package identity

import (
	"time"

	"shanhu.io/std/errcode"
)

func publicKeyValid(k *PublicKey, now time.Time) error {
	if k.NotValidBefore > 0 {
		if now.Before(time.Unix(k.NotValidBefore, 0)) {
			return errcode.InvalidArgf("key not valid yet")
		}
	}
	if now.After(time.Unix(k.NotValidAfter, 0)) {
		return errcode.InvalidArgf("key expired")
	}
	return nil
}
