package signinapi

import (
	"shanhu.io/g/timeutil"
)

// ChallengeRequest is the request to get a challenge.
type ChallengeRequest struct{}

// ChallengeResponse is the response that contains a challenge for the
// client to sign. The challenge normally can only be used once and must be
// used with in a small, limited time window upon issued.
type ChallengeResponse struct {
	Challenge []byte
	Time      *timeutil.Timestamp
}
