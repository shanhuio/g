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
