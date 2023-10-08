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
	"io/fs"
	"os"

	"shanhu.io/g/errcode"
)

type fileStat struct {
	Name         string
	Type         string
	Size         int64
	ModTimestamp int64
	Mode         uint32
	Symlink      string `json:",omitempty"`
}

const (
	fileTypeSrc = "s"
	fileTypeOut = "o"
)

func newOutFileStat(env *env, p string) (*fileStat, error) {
	return newFileStat(env, p, fileTypeOut)
}

func newSrcFileStat(env *env, p string) (*fileStat, error) {
	return newFileStat(env, p, fileTypeSrc)
}

func newFileStat(env *env, p, t string) (*fileStat, error) {
	var f string
	if t == fileTypeOut {
		f = env.out(p)
	} else {
		f = env.src(p)
	}

	info, err := os.Lstat(f) // Does not follow symlink.
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errcode.NotFoundf("%s:%s not found", t, p)
		}
		return nil, err
	}

	var symLink string
	mod := info.Mode()
	if mod&fs.ModeSymlink != 0 { // A sym link.
		dest, err := os.Readlink(f)
		if err != nil {
			return nil, errcode.Annotate(err, "read sym link")
		}
		symLink = dest
	}

	return &fileStat{
		Name:         p,
		Type:         t,
		Size:         info.Size(),
		ModTimestamp: info.ModTime().UnixNano(),
		Mode:         uint32(info.Mode()),
		Symlink:      symLink,
	}, nil
}

func sameFileStat(env *env, stat *fileStat) (bool, error) {
	cur, err := newFileStat(env, stat.Name, stat.Type)
	if err != nil {
		if errcode.IsNotFound(err) {
			return false, nil
		}
		return false, errcode.Annotate(err, "check current")
	}

	same := cur.Size == stat.Size
	same = same && cur.ModTimestamp == stat.ModTimestamp
	same = same && cur.Mode == stat.Mode
	same = same && cur.Symlink == stat.Symlink

	return same, nil
}
