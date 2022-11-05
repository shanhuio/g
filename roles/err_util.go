package roles

import (
	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
)

func altAuthErr(err error, msg string) error {
	if errcode.IsNotFound(err) {
		return errcode.NotFoundf("role not found")
	}
	if errcode.IsUnauthorized(err) || errcode.IsInvalidArg(err) {
		return err
	}
	return aries.AltInternal(err, msg)
}
