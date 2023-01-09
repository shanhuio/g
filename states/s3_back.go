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

package states

import (
	"net/url"

	"shanhu.io/pub/s3util"
)

type s3Back struct {
	client *s3util.Client
}

func newS3Back(client *s3util.Client) *s3Back {
	return &s3Back{client: client}
}

func (b *s3Back) Get(ctx C, key string) ([]byte, error) {
	return b.client.GetBytes(ctx, key)
}

func (b *s3Back) Put(ctx C, key string, data []byte) error {
	return b.client.PutBytes(ctx, key, data)
}

func (b *s3Back) Del(ctx C, key string) error {
	return b.client.Delete(ctx, key)
}

func (b *s3Back) URL() *url.URL {
	return b.client.BaseURL()
}
