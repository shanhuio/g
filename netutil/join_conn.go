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

package netutil

import (
	"context"
	"io"
	"net"
	"sync"
)

// JoinConn joins two net.Conn's. It exits if on any kind of read/write
// error and will close both connections.
func JoinConn(ctx context.Context, c1, c2 net.Conn) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	var closeOnce sync.Once
	closeAll := func() {
		closeOnce.Do(func() {
			c1.Close()
			c2.Close()
		})
	}
	defer closeAll()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	retErr := make(chan error, 3)

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		<-ctx.Done()
		retErr <- ctx.Err()
		closeAll()
	}(ctx)

	var ioWait sync.WaitGroup

	join := func(c1, c2 net.Conn) {
		defer func() {
			closeAll()
			ioWait.Done()
		}()
		if _, err := io.Copy(c1, c2); err != nil {
			retErr <- err
		}
	}

	ioWait.Add(2)
	go join(c1, c2)
	go join(c2, c1)
	ioWait.Wait()

	select {
	case err := <-retErr:
		return err
	default:
		return nil
	}
}
