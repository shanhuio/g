package signinapi

import (
	"time"

	"shanhu.io/g/signer"
	"shanhu.io/g/timeutil"
)

// Request is the request for signing in and creating a session.
type Request struct {
	User        string
	SignedTime  *signer.SignedRSABlock `json:",omitempty"`
	AccessToken string                 `json:",omitempty"`

	TTLDuration *timeutil.Duration `json:",omitempty"`

	TTL int64 `json:",omitempty"` // Nano duration, legacy use.
}

// FillLegacyTTL fills the legacy TTL field, so that it is backwards
// compatible.
func (r *Request) FillLegacyTTL() {
	if r.TTLDuration != nil && r.TTL == 0 {
		r.TTL = r.TTLDuration.Duration().Nanoseconds()
	}
}

// GetTTL gets the TTL. It respects the legacy TTL field.
func (r *Request) GetTTL() time.Duration {
	if r.TTLDuration == nil {
		return time.Duration(r.TTL)
	}
	return timeutil.TimeDuration(r.TTLDuration)
}
