// Copyright (C) 2022  Shanhu Tech Inc.
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
	"shanhu.io/pub/timeutil"
)

// SSHSignInRecord is the record that is being signed
type SSHSignInRecord struct {
	User      string
	Challenge []byte
	TTL       *timeutil.Duration `json:",omitempty"`
}

// SSHSignInRequest is the request to sign in with an SSH certificate
// credential.
type SSHSignInRequest struct {
	RecordBytes []byte // JSON encoded SSHSignInRecord
	Sig         *SSHSignature
	Certificate string
}

// SSHSignature is a copy of *ssh.Signature, it represents an SSH signature.
type SSHSignature struct {
	Format string
	Blob   []byte
	Rest   []byte
}
