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
