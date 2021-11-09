package wsconn

import (
	"net"
	"time"
)

type FrameHandler interface {
	ReadFrame() (b []byte, err error)
	Write(b []byte) (n int, err error)
	Close() error
	ForceClose() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetDeadline(t time.Time) error
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
}
