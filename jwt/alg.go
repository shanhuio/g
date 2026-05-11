package jwt

// Signing algorithm codes.
const (
	AlgHS256 = "HS256" // HMAC + SHA256
	AlgRS256 = "RS256" // RSA + SHA256
)

// The default type string.
const (
	DefaultType = "JWT"
)
