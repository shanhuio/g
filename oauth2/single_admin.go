package oauth2

// SingleAdmin returns a user check function for a single admin.
func SingleAdmin(admin string) func(user string) (any, int, error) {
	return func(user string) (any, int, error) {
		if user == admin {
			return user, 10, nil
		}
		return "", 0, nil
	}
}
