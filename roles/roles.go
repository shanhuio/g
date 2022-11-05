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

package roles

import (
	"sort"
	"time"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/identity"
	"shanhu.io/pub/jwt"
	"shanhu.io/pub/pisces"
	"shanhu.io/pub/rand"
	"shanhu.io/pub/roles/rolesapi"
	"shanhu.io/pub/timeutil"
)

// Roles contains tables that saves user's identity.
// It provides a service for users to
//   - create a role
//   - create a one-time code to setup credential keys
//   - register a new public key
//   - check a self-signed ID token.
type Roles struct {
	t *pisces.KV

	host string

	passCodeExpiry time.Duration
}

// New creates the roles table.
func New(b *pisces.Tables) *Roles {
	return NewWithName(b, "roles")
}

// NewWithName creates the roles table using the given table name.
func NewWithName(b *pisces.Tables, name string) *Roles {
	return &Roles{
		t:              b.NewKV(name),
		passCodeExpiry: 10 * time.Minute,
	}
}

type role struct {
	Role *rolesapi.Role

	Identity *identity.Identity `json:",omitempty"`

	// PassCode is a one-time passcode that can be used for
	// registering the identity.
	PassCode *passCode `json:",omitempty"`
}

// New creates a new empty role. Its identity is empty.
func (b *Roles) New(name string, t time.Time) error {
	pub := &rolesapi.Role{
		Name:       name,
		TimeCreate: timeutil.NewTimestamp(t),
	}
	r := &role{Role: pub}
	return b.t.Add(name, r)
}

func (b *Roles) get(name string) (*role, error) {
	r := new(role)
	if err := b.t.Get(name, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Get returns the role's info.
func (b *Roles) Get(name string) (*rolesapi.Role, error) {
	r, err := b.get(name)
	if err != nil {
		return nil, err
	}
	return r.Role, nil
}

// Remove removes a role.
func (b *Roles) Remove(name string) error { return b.t.Remove(name) }

func (b *Roles) setDisabled(name string, v bool) error {
	return b.mutate(name, func(r *role) error {
		r.Role.Disabled = v
		return nil
	})
}

// Disable disables a role.
func (b *Roles) Disable(name string) error {
	return b.setDisabled(name, true)
}

// Enable enables a role.
func (b *Roles) Enable(name string) error {
	return b.setDisabled(name, false)
}

// List lists all roles.
func (b *Roles) List() ([]*rolesapi.Role, error) {
	items := make([]*rolesapi.Role, 0)
	it := &pisces.Iter{
		Make: func() interface{} { return new(role) },
		Do: func(_ string, v interface{}) error {
			item := v.(*role).Role
			if item == nil {
				item = new(rolesapi.Role)
			}
			items = append(items, item)
			return nil
		},
	}
	if err := b.t.Walk(it); err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		sec1 := items[i].TimeCreate.Sec
		sec2 := items[j].TimeCreate.Sec
		if sec1 == sec2 {
			return items[i].Name < items[j].Name
		}
		return sec1 < sec2
	})
	return items, nil
}

// SetPassCodeExpiry sets the passcode expiry.
func (b *Roles) SetPassCodeExpiry(d time.Duration) {
	const zero = time.Duration(0)
	if d < zero {
		d = zero
	}
	b.passCodeExpiry = d
}

// SetHostDomain sets the host domain for checking JWT.
func (b *Roles) SetHostDomain(domain string) { b.host = domain }

func (b *Roles) mutate(name string, f func(r *role) error) error {
	r := new(role)
	return b.t.Mutate(name, r, func(v interface{}) error {
		return f(v.(*role))
	})
}

// NewPassCode creates a new passcode for roles, with the given timestamp.
func (b *Roles) NewPassCode(name string, now time.Time) (
	*rolesapi.PassCode, error,
) {
	const buffer = 1 * time.Minute
	code := &passCode{
		Code:   rand.Digits(8),
		Valid:  timeutil.NewTimestamp(now.Add(-buffer)),
		Expire: timeutil.NewTimestamp(now.Add(b.passCodeExpiry)),
	}
	if err := b.mutate(name, func(r *role) error {
		if r.Role.Disabled {
			return errcode.InvalidArgf("role is disabled")
		}
		r.PassCode = code
		return nil
	}); err != nil {
		return nil, err
	}
	return code.public(), nil
}

// GetPassCode fetches the role's current pass code if any.
func (b *Roles) GetPassCode(name string) (*rolesapi.PassCode, error) {
	r, err := b.get(name)
	if err != nil {
		return nil, err
	}
	if r.PassCode == nil {
		return nil, nil
	}
	return r.PassCode.public(), nil
}

// SetupWithCode sets up an identity with the given pass code.
// The identity is set only when the given pass code is valid.
func (b *Roles) SetupWithCode(
	name string, id *identity.Identity, code string, t time.Time,
) error {
	return b.mutate(name, func(r *role) error {
		if r.Role.Disabled {
			return errcode.InvalidArgf("role is disabled")
		}
		if r.PassCode != nil {
			r.PassCode.Tried++
		}
		if err := checkPassCode(code, r.PassCode, t); err != nil {
			return err
		}
		r.Identity = id
		r.PassCode.Consumed = true
		return nil
	})
}

// VerifySelfToken checks the self-signed JWT token.
func (b *Roles) VerifySelfToken(name, token string, t time.Time) (
	*jwt.Token, error,
) {
	r, err := b.get(name)
	if err != nil {
		return nil, err
	}
	if r.Role.Disabled {
		return nil, errcode.Unauthorizedf("role is disabled")
	}
	if r.Identity == nil {
		return nil, errcode.Unauthorizedf("no identity found")
	}
	return identity.VerifySelfToken(token, name, b.host, r.Identity, t)
}
