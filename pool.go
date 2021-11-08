package wsconn

import (
	"bytes"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func NewWSConn(conn *websocket.Conn) net.Conn {
	wc := &wsconn{
		conn:  conn,
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

type wsconn struct {
	conn *websocket.Conn
	mu sync.RWMutex
	rmu sync.Mutex
	wmu sync.Mutex
	buf bytes.Buffer
	come chan bool
	done chan struct{}
	errCh chan error
}

func (w *wsconn) readConn() {
	for {
		w.rmu.Lock()

		t, b, err := w.conn.ReadMessage()
		if err != nil {
			w.errCh <- err
			w.rmu.Unlock()
			return
		}

		switch t {
		case websocket.BinaryMessage, websocket.TextMessage:
			w.buf.Write(b)
			select {
			case w.come <- true:
			default:
			}
		case websocket.CloseMessage:
			w.rmu.Unlock()
			_ = w.Close()
			return
		case websocket.PingMessage, websocket.PongMessage:
		}

		w.rmu.Unlock()
	}
}

func (w *wsconn) Read(b []byte) (n int, err error) {
	for {
		w.mu.RLock()
		w.rmu.Lock()
		if w.buf.Len() != 0 {
			break
		}
		w.rmu.Unlock()
		w.mu.RUnlock()
		select {
		case <- w.done:
			return 0, io.EOF
		case <- w.come:
			continue
		case e := <- w.errCh:
			return 0, e
		}
	}
	defer w.mu.RUnlock()
	defer w.rmu.Unlock()

	return w.buf.Read(b)
}

func (w *wsconn) Write(b []byte) (n int, err error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	w.wmu.Lock()
	defer w.wmu.Unlock()

	err = w.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *wsconn) Close() error {
	log.Println("Getting Lock")
	w.mu.Lock()
	log.Println("Got it!")
	defer w.mu.Unlock()
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		return err
	}
	close(w.done)
	return w.conn.Close()
}

func (w *wsconn) ForceClose() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	close(w.done)
	return w.conn.Close()
}

func (w *wsconn) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *wsconn) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *wsconn) SetDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	err := w.conn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return w.conn.SetWriteDeadline(t)
}

func (w *wsconn) SetReadDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.SetReadDeadline(t)
}

func (w *wsconn) SetWriteDeadline(t time.Time) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.conn.SetWriteDeadline(t)
}
