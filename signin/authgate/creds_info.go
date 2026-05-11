package authgate

import (
	"shanhu.io/g/aries"
)

// CredsInfo is the user credential information got from gate checking.
type CredsInfo struct {
	Valid       bool
	NeedRefresh bool

	TokenType string
	User      string
	UserLevel int

	Data interface{}
}

// ApplyCredsInfo applies the credential into the aries context.
func ApplyCredsInfo(c *aries.C, info *CredsInfo) {
	if !info.Valid {
		c.User = ""
		c.UserLevel = 0
		return
	}

	c.User = info.User
	c.UserLevel = info.UserLevel
	if info.Data != nil {
		c.UserData = info.Data
	}
}
