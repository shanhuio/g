package rolesapi

import (
	"shanhu.io/g/timeutil"
)

// Role contains the roles meta information.
type Role struct {
	Name       string
	TimeCreate *timeutil.Timestamp
	Disabled   bool `json:",omitempty"`
}

// PassCode contains the info of a pass code.
type PassCode struct {
	Code   string
	Valid  *timeutil.Timestamp
	Expire *timeutil.Timestamp

	TriedTooManyTimes bool
}
