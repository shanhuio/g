package pisces

import (
	"errors"

	"shanhu.io/g/errcode"
)

// ErrCancel is the error when the operation is canclled.
var ErrCancel = errors.New("operation cancelled")

// ErrUnordered is the error when an ordered table operation is
// operated on an unordered table.
var ErrUnordered = errors.New("the index is unordered")

var (
	notFound      = errcode.NotFoundf("not found")
	multiAffected = errcode.Internalf("multiple entries affected")
)
