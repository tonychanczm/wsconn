package wsconn

import (
	"bytes"
	"io"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func NewWSConn(conn *websocket.Conn) net.Conn {
	return NewConnectAdapter(NewWSFrameHandler(conn))
}

func NewConnectAdapter(fh FrameHandler) net.Conn {
	wc := &connectAdapter{
		fh:    fh,
		mu:    sync.RWMutex{},
		rmu:   sync.Mutex{},
		wmu:   sync.Mutex{},
		buf:   bytes.Buffer{},
		come:  make(chan bool, 1),
		done:  make(chan struct{}, 1),
		errCh: make(chan error, 1),
	}
	go wc.readConn()
	return wc
}

type connectAdapter struct {
	fh        FrameHandler
	mu        sync.RWMutex
	rmu       sync.Mutex
	wmu       sync.Mutex
	buf       bytes.Buffer
	come      chan bool
	done      chan struct{}
	onceClose sync.Once
	errCh     chan error
}

func (w *connectAdapter) close() {
	w.onceClose.Do(func() {
		close(w.done)
	})
}

func (w *connectAdapter) readConn() {
	for {
		b, err := w.fh.ReadFrame()
		if err != nil {
			if err == io.EOF {
				w.close()
				return
			}
			w.errCh <- err
			return
		}

		w.rmu.Lock()
		w.buf.Write(b)
		w.rmu.Unlock()
		select {
		case w.come <- true:
		default:
		}

	}
}

func (w *connectAdapter) Read(b []byte) (n int, err error) {
	for {
		w.mu.RLock()
		w.rmu.Lock()
		if w.buf.Len() != 0 {
			break
		}
		w.rmu.Unlock()
		w.mu.RUnlock()
		select {
		case <-w.done:
			return 0, io.EOF
		case <-w.come:
			continue
		case e := <-w.errCh:
			return 0, e
		}
	}
	defer w.mu.RUnlock()
	defer w.rmu.Unlock()
	return w.buf.Read(b)
}

func (w *connectAdapter) Write(b []byte) (n int, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	w.wmu.Lock()
	defer w.wmu.Unlock()

	return w.fh.Write(b)
}

func (w *connectAdapter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.close()
	return w.fh.Close()
}

func (w *connectAdapter) ForceClose() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	defer w.close()
	return w.fh.ForceClose()
}

func (w *connectAdapter) LocalAddr() net.Addr {
	return w.fh.LocalAddr()
}

func (w *connectAdapter) RemoteAddr() net.Addr {
	return w.fh.RemoteAddr()
}

func (w *connectAdapter) SetDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	err := w.fh.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return w.fh.SetWriteDeadline(t)
}

func (w *connectAdapter) SetReadDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fh.SetReadDeadline(t)
}

func (w *connectAdapter) SetWriteDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.fh.SetWriteDeadline(t)
}
