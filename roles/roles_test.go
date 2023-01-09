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

package roles

import (
	"testing"

	"context"
	"time"

	"shanhu.io/pub/identity"
	"shanhu.io/pub/pisces"
	"shanhu.io/pub/timeutil"
)

func TestRoles(t *testing.T) {
	tables := pisces.NewMemTables()
	b := New(tables)

	const name = "h8liu"
	const host = "shanhu.io"
	b.SetHostDomain(host)

	now := time.Now()
	nowFunc := func() time.Time { return now }

	if err := b.New(name, now); err != nil {
		t.Fatal("create role: ", err)
	}

	r, err := b.Get(name)
	if err != nil {
		t.Fatal("get role: ", err)
	}
	if r.Name != name {
		t.Errorf("got role %q, want %q", r.Name, name)
	}
	if timeCreate := timeutil.Time(r.TimeCreate); !timeCreate.Equal(now) {
		t.Errorf("create time, got %q, want %q", timeCreate, now)
	}

	code, err := b.NewPassCode(name, now)
	if err != nil {
		t.Fatal("get new passcode: ", err)
	}

	core := identity.NewMemCore(nowFunc)
	id, err := core.Init(identity.SingleKeyCoreConfig(now.Add(time.Hour)))
	if err != nil {
		t.Fatal("init id: ", err)
	}

	if err := b.SetupWithCode(name, id, code.Code, now); err != nil {
		t.Fatal("setup id: ", err)
	}

	ctx := context.Background()
	self, err := identity.SignSelf(ctx, core, name, host, now)
	if err != nil {
		t.Fatal("self-sign ID token: ", err)
	}

	if _, err := b.VerifySelfToken(ctx, name, self, now); err != nil {
		t.Fatal("verify self-sign ID token: ", err)
	}
}
