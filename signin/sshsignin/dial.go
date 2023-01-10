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

package sshsignin

import (
	"context"
	"encoding/json"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/signin/signinapi"
	"shanhu.io/pub/strutil"
)

// Config contains the configuration to sign in with an SSH
// certificate.
type Config struct {
	User string // Default using SHANHU_USER or system user name.

	Agent      agent.ExtendedAgent // Default using SSH_AUTH_SOCK
	KeyComment string              // Default is "shanhu"
}

func (c *Config) user() (string, error) {
	if c.User != "" {
		return c.User, nil
	}
	return SysUser()
}

func (c *Config) agent() (agent.ExtendedAgent, error) {
	if c.Agent != nil {
		return c.Agent, nil
	}
	return SysAgent()
}

// FindKey finds the agent key in the SSH agent that has the exact comment.
func FindKey(ag agent.Agent, comment string) (*agent.Key, error) {
	keys, err := ag.List()
	if err != nil {
		return nil, errcode.Annotate(err, "list keys")
	}
	for _, k := range keys {
		if k.Comment == comment {
			return k, nil
		}
	}
	return nil, errcode.NotFoundf("%q not found", comment)
}

func findKey(ag agent.Agent, comment string) (*agent.Key, error) {
	k, err := FindKey(ag, comment)
	if err != nil {
		if errcode.IsNotFound(err) { // converts not found error into internal.
			return nil, errcode.Add(errcode.Internal, err)
		}
		return nil, err
	}
	return k, nil
}

// Dial signs in a server and returns the credentials.
func Dial(ctx context.Context, server string, config *Config) (
	*httputil.Client, error,
) {
	user, err := config.user()
	if err != nil {
		return nil, errcode.Annotate(err, "get user name")
	}

	ag, err := config.agent()
	if err != nil {
		return nil, errcode.Annotate(err, "get SSH agent")
	}

	keyComment := strutil.Default(config.KeyComment, "shanhu")
	key, err := findKey(ag, keyComment)
	if err != nil {
		return nil, errcode.Annotate(err, "find key")
	}
	if t := key.Type(); t != ssh.CertAlgoRSAv01 {
		return nil, errcode.Internalf("unexpected key type %q", t)
	}

	client, err := httputil.NewClient(server)
	if err != nil {
		return nil, errcode.Annotate(err, "make http client")
	}

	chReq := &signinapi.ChallengeRequest{}
	chResp := new(signinapi.ChallengeResponse)
	const chPath = "/ssh/challenge"
	if err := client.CallContext(ctx, chPath, chReq, chResp); err != nil {
		return nil, errcode.Annotate(err, "get challenge")
	}

	record := &signinapi.SSHSignInRecord{
		User:      user,
		Challenge: chResp.Challenge,
	}
	recordBytes, err := json.Marshal(record)
	if err != nil {
		return nil, errcode.Annotate(err, "marshal signin record")
	}
	const signFlag = agent.SignatureFlagRsaSha256
	sig, err := ag.SignWithFlags(key, recordBytes, signFlag)
	if err != nil {
		return nil, errcode.Annotate(err, "sign signin record")
	}

	req := signinapi.SSHSignInRequest{
		RecordBytes: recordBytes,
		Sig: &signinapi.SSHSignature{
			Format: sig.Format,
			Blob:   sig.Blob,
			Rest:   sig.Rest,
		},
		Certificate: string(ssh.MarshalAuthorizedKey(key)),
	}

	creds := new(signinapi.Creds)
	if err := client.CallContext(ctx, "/ssh/signin", req, creds); err != nil {
		return nil, err
	}

	client.TokenSource = httputil.NewStaticToken(creds.Token)
	return client, nil
}
