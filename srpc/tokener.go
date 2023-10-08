// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package srpc

import (
	"context"
	"time"

	"shanhu.io/g/timeutil"
)

// Tokener in an interface that issues token for authentication.
type Tokener interface {
	// Token refreshes the token. If lastToken is empty, it will issue a new
	// token. Otherwise, it will refresh the token using the lastToken.
	Token(ctx context.Context) (string, error)
}

// TimedToken is a token that has a refresh time and an expire time.
type TimedToken struct {
	Token string

	// Refresh is the time when the token needs refreshing.
	// This does not necessarily mean that the token is expired.
	// If Refresh is nil, it means the token always needs refreshing.
	Refresh *timeutil.Timestamp
}

// TimedTokener is an interface that issues token for authentication.
// The token has a refresh time and an expire time.
type TimedTokener interface {
	TimedToken(ctx context.Context) (*TimedToken, error)
}

// CachingTokener is a tokener that caches the token based on
// the refresh time and expire time.
type CachingTokener struct {
	now     func() time.Time
	cache   *TimedToken
	tokener TimedTokener
}

// NewCachedTokener returns a caching tokener.
func NewCachedTokener(tt TimedTokener) *CachingTokener {
	return &CachingTokener{
		now:     time.Now,
		tokener: tt,
	}
}

// SetClock sets a custom clock function for the tokener. time.Now will be used.
// if not set, or set to nil.
func (t *CachingTokener) SetClock(f func() time.Time) {
	t.now = timeutil.NowFunc(f)
}

// Token returns a token. If there is a cached token that does not yet need
// refreshing, it will return the cached one. Otherwise, it will refresh the
// token and return the new one.
func (t *CachingTokener) Token(ctx context.Context) (string, error) {
	if t.cache != nil && t.cache.Refresh != nil {
		now := t.now()
		if now.Before(t.cache.Refresh.Time()) {
			return t.cache.Token, nil
		}
	}

	tt, err := t.tokener.TimedToken(ctx)
	if err != nil {
		return "", err
	}

	if t.cache.Refresh != nil {
		t.cache = tt
	}
	return t.cache.Token, nil
}
