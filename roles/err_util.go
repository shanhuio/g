package roles

import (
	"shanhu.io/g/aries"
	"shanhu.io/g/errcode"
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
