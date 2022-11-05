package roles

import (
	"crypto/subtle"
	"time"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/roles/rolesapi"
	"shanhu.io/pub/timeutil"
)

const passCodeMaxTries = 10

type passCode struct {
	Code     string
	Valid    *timeutil.Timestamp
	Expire   *timeutil.Timestamp
	Consumed bool `json:",omitempty"`
	Tried    int  `json:",omitempty"`
}

func (c *passCode) public() *rolesapi.PassCode {
	if c.Consumed {
		return nil
	}
	return &rolesapi.PassCode{
		Code:              c.Code,
		Valid:             timeutil.CopyTimestamp(c.Valid),
		Expire:            timeutil.CopyTimestamp(c.Expire),
		TriedTooManyTimes: c.Tried > passCodeMaxTries,
	}
}

func subtleStringEq(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) != 0
}

func checkPassCode(
	claim string, code *passCode, now time.Time,
) error {
	if claim == "" {
		return errcode.Unauthorizedf("empty passcode")
	}
	if code == nil {
		return errcode.Unauthorizedf("no passcode set")
	}
	if code.Valid == nil {
		return errcode.Internalf("passcode valid time missing")
	}
	if code.Expire == nil {
		return errcode.Internalf("passcode expire time missing")
	}

	if code.Tried > passCodeMaxTries {
		return errcode.Unauthorizedf("passcode wrong too many times")
	}
	if code.Consumed {
		return errcode.Unauthorizedf("passcode already consumed")
	}

	valid := code.Valid.Time()
	if now.Before(valid) {
		return errcode.Unauthorizedf("passcode not valid yet")
	}
	expire := code.Expire.Time()
	if now.After(expire) {
		return errcode.Unauthorizedf("passcode expired")
	}

	if !subtleStringEq(code.Code, claim) {
		return errcode.Unauthorizedf("passcode incorrect")
	}
	return nil
}
