package netutil

import (
	"testing"

	"bytes"
	"context"
	"crypto/rand"
	"io"
	"net"

	"shanhu.io/g/errcode"
)

func serveEchoBack(c net.Conn, messageLen int) error {
	recv := make([]byte, messageLen)
	n, err := io.ReadFull(c, recv)
	if err != nil {
		return errcode.Annotate(err, "read b1")
	}
	if n != messageLen {
		return errcode.Internalf("read message got %d bytes", n)
	}
	if _, err := c.Write(recv); err != nil {
		return errcode.Annotate(err, "write back to b1")
	}
	return nil
}

func TestJoinConn(t *testing.T) {
	a1, a2 := net.Pipe()
	b2, b1 := net.Pipe()

	ctx := context.Background()

	joinConnErr := make(chan error)
	go func(ctx context.Context) {
		joinConnErr <- JoinConn(ctx, a2, b2)
	}(ctx)

	const messageLen = 16

	serveErr := make(chan error)
	go func() {
		serveErr <- serveEchoBack(b1, messageLen)
	}()

	msg := make([]byte, messageLen)
	if _, err := rand.Read(msg); err != nil {
		t.Fatal("prepare message: ", err)
	}

	if _, err := a1.Write(msg); err != nil {
		t.Fatal("write message to a1: ", err)
	}

	recv := make([]byte, messageLen)
	n, err := io.ReadFull(a1, recv)
	if err != nil {
		t.Fatal("read message: ", err)
	}
	if n != messageLen {
		t.Fatalf("read back got %d bytes", n)
	}

	if !bytes.Equal(recv, msg) {
		t.Errorf("sent %x, recieved %x", recv, msg)
	}

	if err := a1.Close(); err != nil {
		t.Error("close conn: ", err)
	}

	if err := <-serveErr; err != nil {
		t.Error("serve error: ", err)
	}
	if err := <-joinConnErr; err != nil {
		if err != io.ErrClosedPipe {
			t.Error("join conn error: ", err)
		}
	}
}
