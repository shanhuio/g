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

package authgate

import (
	"crypto/rand"
	"io"
	"time"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/timeutil"
)

// ChallengerConfig is the configuration to create a challenge source.
type ChallengerConfig struct {
	Signer *signer.Signer
	Now    func() time.Time
	Rand   io.Reader
	Window time.Duration
}

// Challenger is a source that can serve challenges.
type Challenger struct {
	signer *signer.Signer
	now    func() time.Time
	rand   io.Reader
	window time.Duration
}

// NewChallenger creates a challenge source.
func NewChallenger(config *ChallengerConfig) *Challenger {
	r := config.Rand
	if r == nil {
		r = rand.Reader
	}
	w := config.Window
	if w <= time.Duration(0) {
		w = 30 * time.Second
	}
	return &Challenger{
		signer: config.Signer,
		now:    timeutil.NowFunc(config.Now),
		rand:   r,
		window: w,
	}
}

// Serve serves a challenge.
func (s *Challenger) Serve(
	c *aries.C, req *signinapi.ChallengeRequest,
) (*signinapi.ChallengeResponse, error) {
	t := s.now()
	signed, ch, err := s.signer.NewSignedChallenge(t, s.rand)
	if err != nil {
		return nil, err
	}
	return &signinapi.ChallengeResponse{
		Challenge: signed,
		Time:      ch.T,
	}, nil
}

// Check checks if a challenge is valid.
func (s *Challenger) Check(bs []byte) (*timeutil.Challenge, error) {
	now := s.now()
	return s.signer.CheckChallenge(bs, now, s.window)
}
