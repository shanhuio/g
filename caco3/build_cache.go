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

package caco3

import (
	"time"

	"shanhu.io/g/errcode"
	"shanhu.io/g/pisces"
	"shanhu.io/g/timeutil"

	_ "modernc.org/sqlite" // sqlite db driver
)

type buildCache struct {
	tables *pisces.Tables
	cache  *pisces.KV
	expire time.Duration
	clock  func() time.Time
}

func newBuildCache(f string) (*buildCache, error) {
	tables, err := pisces.OpenSqlite3Tables(f)
	if err != nil {
		return nil, errcode.Annotate(err, "open cache table")
	}

	cache := tables.NewKV("build_cache")
	if err := tables.CreateMissing(); err != nil {
		return nil, errcode.Annotate(err, "create cache tables")
	}

	return &buildCache{
		expire: time.Hour * 24 * 7,
		tables: tables,
		cache:  cache,
	}, nil
}

type buildCacheEntry struct {
	Key        string              `json:"K"`
	CreateTime *timeutil.Timestamp `json:"T"`
	Built      *built              `json:"B"`
}

func (c *buildCache) put(k string, out *built) error {
	t := timeutil.ReadTime(c.clock)
	entry := &buildCacheEntry{
		Key:        k,
		Built:      out,
		CreateTime: timeutil.NewTimestamp(t),
	}
	return c.cache.Replace(k, entry)
}

var errNotFoundInCache = errcode.NotFoundf("not found in cache")

func (c *buildCache) get(k string) (*built, error) {
	entry := new(buildCacheEntry)
	if err := c.cache.Get(k, entry); err != nil {
		if errcode.IsNotFound(err) {
			return nil, errNotFoundInCache
		}
		return nil, errcode.Annotate(err, "get from cache")
	}

	now := timeutil.ReadTime(c.clock)
	expire := timeutil.Time(entry.CreateTime).Add(c.expire)
	if now.Before(expire) {
		return entry.Built, nil
	}
	return nil, errNotFoundInCache
}

func (c *buildCache) remove(k string) error {
	if err := c.cache.Remove(k); err != nil {
		if errcode.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}
