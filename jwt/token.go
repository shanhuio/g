package jwt

// Token is a parsed JWT token.
type Token struct {
	Header    *Header
	ClaimSet  *ClaimSet
	Payload   []byte
	Signature []byte
}
