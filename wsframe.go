package wsconn

import (
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

func NewWSFrameHandler(conn *websocket.Conn) FrameHandler {
	return &wsFrame{conn: conn}
}

type wsFrame struct {
	conn *websocket.Conn
}

func (w *wsFrame) ReadFrame() (b []byte, err error) {
	for {
		t, b, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				return nil, io.EOF
			}
			return nil, err
		}

		switch t {
		case websocket.BinaryMessage, websocket.TextMessage:
			return b, nil
		case websocket.CloseMessage:
			_ = w.Close()
			return nil, io.EOF
		case websocket.PingMessage, websocket.PongMessage:
		}
	}
}

func (w *wsFrame) Write(b []byte) (n int, err error) {
	err = w.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (w *wsFrame) Close() error {
	err := w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		if err == websocket.ErrCloseSent {
			return nil
		}
		return err
	}
	return w.conn.Close()
}

func (w *wsFrame) ForceClose() error {
	_ = w.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	return w.conn.Close()
}

func (w *wsFrame) LocalAddr() net.Addr {
	return w.conn.LocalAddr()
}

func (w *wsFrame) RemoteAddr() net.Addr {
	return w.conn.RemoteAddr()
}

func (w *wsFrame) SetDeadline(t time.Time) error {
	err := w.conn.SetReadDeadline(t)
	if err != nil {
		return err
	}
	return w.conn.SetWriteDeadline(t)
}

func (w *wsFrame) SetReadDeadline(t time.Time) error {
	return w.conn.SetReadDeadline(t)
}

func (w *wsFrame) SetWriteDeadline(t time.Time) error {
	return w.conn.SetWriteDeadline(t)
}
