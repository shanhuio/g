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
)

// Token is a bearer token. It also embeds a refresh recommendation
// and an expiration time point.
type Token struct {
	Token string

	// If the current time is after this time point, the caller should try to
	// refresh the token using the refresh function.
	Refresh time.Time

	// If time is after this time point, the caller should stop using this
	// token for requests.
	Expire time.Time
}

// Tokener in an interface that issues token for authentication.
type Tokener interface {
	// Token refreshes the token. If lastToken is empty, it will issue a new
	// token. Otherwise, it will refresh the token using the lastToken.
	Token(ctx context.Context, lastToken string) (*Token, error)
}
