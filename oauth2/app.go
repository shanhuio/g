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
