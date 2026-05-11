package oauth2

// App stores the configuration of a general oauth2 application.
type App struct {
	ID          string
	Secret      string
	RedirectURL string `json:",omitempty"`

	Scopes []string `json:",omitempty"`

	// Used only in GitHub OAuth2
	WithEmail bool `json:",omitempty"`

	// Used only in Google OAuth2
	WithProfile bool `json:",omitempty"`
}
