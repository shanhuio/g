package signinapi

import (
	"shanhu.io/g/timeutil"
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
