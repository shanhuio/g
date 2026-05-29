package jwt

import (
	"time"

	"shanhu.io/std/errcode"
)

// CheckTime checks if the token's claims is in valid at time now.
func CheckTime(claims *ClaimSet, now time.Time) (time.Duration, error) {
	// Issued time must be after now.
	// In case of small clock error, gives a 5 minute grace period.
	issued := time.Unix(claims.Iat, 0).Add(-5 * time.Minute)
	if !issued.Before(now) {
		return 0, errcode.Unauthorizedf("token issued in the future")
	}

	expires := time.Unix(claims.Exp, 0)
	if now.After(expires) {
		return 0, errcode.Unauthorizedf("token expired")
	}
	return expires.Sub(now), nil
}
