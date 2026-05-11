package identity

import (
	"fmt"
)

// UserAtDomain returns the string of user@domain.
func UserAtDomain(user, domain string) string {
	if domain == "" {
		return user
	}
	return fmt.Sprintf("%s@%s", user, domain)
}
