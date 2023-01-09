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

package signin

import (
	"crypto/rand"
	"io"
	"time"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/timeutil"
)

// ChallengeSourceConfig is the configuration to create a challenge source.
type ChallengeSourceConfig struct {
	Signer *signer.Signer
	Now    func() time.Time
	Rand   io.Reader
}

// Challenger is a source that can serve challenges.
type Challenger struct {
	signer  *signer.Signer
	nowFunc func() time.Time
	rand    io.Reader
}

// NewChallenger creates a challenge source.
func NewChallenger(config *ChallengeSourceConfig) *Challenger {
	r := config.Rand
	if r == nil {
		r = rand.Reader
	}
	return &Challenger{
		signer:  config.Signer,
		nowFunc: timeutil.NowFunc(config.Now),
		rand:    r,
	}
}

// Serve serves a challenge.
func (s *Challenger) Serve(
	c *aries.C, req *signinapi.ChallengeRequest,
) (*signinapi.ChallengeResponse, error) {
	t := s.nowFunc()
	signed, ch, err := s.signer.NewSignedChallenge(t, s.rand)
	if err != nil {
		return nil, err
	}
	return &signinapi.ChallengeResponse{
		Challenge: signed,
		Time:      ch.T,
	}, nil
}
