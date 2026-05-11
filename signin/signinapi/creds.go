package signinapi

import (
	"time"

	"shanhu.io/g/timeutil"
)

// Creds is the response for signing in. It saves the user credentials.
type Creds struct {
	User        string
	Token       string
	ExpiresTime *timeutil.Timestamp `json:",omitempty"`

	Expires int64 `json:",omitempty"` // Nanosecond timestamp, legacy use.
}

// FixTime fixes timestamps.
func (c *Creds) FixTime() {
	if c.ExpiresTime == nil && c.Expires != 0 {
		t := time.Unix(0, c.Expires)
		c.ExpiresTime = timeutil.NewTimestamp(t)
	}
}
