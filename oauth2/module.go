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

package oauth2

import (
	"log"
	"path"
	"sort"
	"time"

	"shanhu.io/g/aries"
	"shanhu.io/g/errcode"
	"shanhu.io/g/signer"
	"shanhu.io/g/signin"
	"shanhu.io/g/signin/authgate"
)

// Module is a module that handles stuff related to oauth.
type Module struct {
	config    *Config
	gate      *authgate.Gate
	providers []provider
	pubKey    *authgate.LegacyExchange

	redirect       string
	signInRedirect string

	clients map[string]*Client

	router *aries.Router
}

// NewModule creates a new oauth module with the given config.
func NewModule(config *Config) *Module {
	redirect := config.Redirect
	if redirect == "" {
		redirect = "/"
	}
	signInRedirect := config.SignInRedirect
	if signInRedirect == "" {
		signInRedirect = redirect
	}

	gate := authgate.New(&authgate.Config{
		SessionKey:      config.SessionKey,
		SessionLifeTime: config.SessionLifeTime,
		Check:           config.Check,
	})

	ret := &Module{
		config:         config,
		gate:           gate,
		redirect:       redirect,
		signInRedirect: signInRedirect,
		clients:        make(map[string]*Client),
	}

	if config.KeyRegistry != nil {
		ret.pubKey = authgate.NewLegacyExchange(
			gate, config.KeyRegistry,
		)
	}

	const ttl time.Duration = time.Hour
	states := signer.NewSessions(config.StateKey, ttl)

	addProvider := func(p provider) {
		ret.providers = append(ret.providers, p)
	}
	if config.GitHub != nil {
		addProvider(newGitHub(config.GitHub, states))
	}
	if config.Google != nil {
		addProvider(newGoogle(config.Google, states))
	}
	if config.DigitalOcean != nil {
		addProvider(newDigitalOcean(config.DigitalOcean, states))
	}

	ret.router = ret.makeRouter()

	return ret
}

// Serve serves the routes for signing in and callbacks.
func (m *Module) Serve(c *aries.C) error { return m.router.Serve(c) }

// Methods returns the list of supported methods.
func (m *Module) Methods() []string {
	var ms []string
	for k := range m.clients {
		ms = append(ms, k)
	}
	sort.Strings(ms)
	return ms
}

func (m *Module) addProvider(r *aries.Router, p provider) {
	c := p.client()
	method := c.Method()
	m.clients[method] = c
	signIn := newSignInHandler(c, m.signInRedirect)
	r.Get(path.Join(method, "signin"), signIn.Serve)
	r.Get(path.Join(method, "callback"), m.callbackHandler(method, p))
}

func (m *Module) signOut(c *aries.C) error {
	if pre := m.config.PreSignOut; pre != nil {
		if err := pre(c); err != nil {
			return err
		}
	}
	authgate.ClearCookie(c)
	c.Redirect(m.redirect)
	return nil
}

func (m *Module) makeRouter() *aries.Router {
	r := aries.NewRouter()
	r.Get("signout", m.signOut)
	if bypass := m.config.Bypass; bypass != "" {
		r.Get("signin-bypass", func(c *aries.C) error {
			m.gate.SetupCookie(c, bypass)
			c.Redirect(m.signInRedirect)
			return nil
		})
	}
	if m.pubKey != nil {
		r.Call("pubkey/signin", m.pubKey.Exchange)
	}
	for _, p := range m.providers {
		m.addProvider(r, p)
	}
	return r
}

// Auth makes a aries.Auth that executes the oauth flow on the server side.
func (m *Module) Auth() aries.Auth { return m }

func (m *Module) signInCheck(
	c *aries.C, u *UserMeta, purpose string,
) (string, error) {
	if f := m.config.SignInCheck; f != nil {
		return f(c, u, purpose)
	}
	return u.ID, nil // default login check allows everyone.
}

func (m *Module) signIn(c *aries.C, user *UserMeta, state *State) error {
	id, err := m.signInCheck(c, user, state.Purpose)
	if err != nil {
		return err
	}
	if id == "" {
		return nil
	}
	if !state.NoCookie {
		m.gate.SetupCookie(c, id)
	}
	c.Redirect(state.Dest)
	return nil
}

// Token returns a new session token for user that expires in ttl.
func (m *Module) Token(user string, ttl time.Duration) *signin.Token {
	return m.gate.Token(user, ttl)
}

// Setup sets up the credentials for the request.
func (m *Module) Setup(c *aries.C) error { return m.gate.Setup(c) }

// SetupCookie sets up the session gate's cookie.
func (m *Module) SetupCookie(c *aries.C, user string) {
	m.gate.SetupCookie(c, user)
}

// SignIn redirects the incoming request to a particular client's sign-in
// URL. If the client is not found, it redirects to default redirect page.
func (m *Module) SignIn(c *aries.C, method string, s *State) error {
	client, ok := m.clients[method]
	if !ok {
		c.Redirect(m.redirect)
		return nil
	}
	c.Redirect(client.SignInURL(s))
	return nil
}

func (m *Module) callbackHandler(method string, x metaExchange) aries.Func {
	return func(c *aries.C) error {
		user, state, err := x.callback(c)
		if err != nil {
			log.Printf("%s callback: %s", method, err)
			return errcode.Internalf("%s callback failed", method)
		}
		if user == nil {
			return errcode.Internalf(
				"%s callback: get user info failed", method,
			)
		}
		return m.signIn(c, user, state)
	}
}
