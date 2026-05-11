package aries

// Auth is an authentication service that sets up the authentication
// context before serving.
type Auth interface {
	Service

	// Setup sets up the authentication in context.
	Setup(c *C) error
}
