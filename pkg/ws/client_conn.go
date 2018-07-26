package ws

import (
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/websocket"
	"github.com/linlycode/olcode/pkg/hubpkg"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512 * 1024
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// AsyncHandler handles a ws connection async
type AsyncHandler interface {
	AsyncRun()
}

type mHandler interface {
	Handle(msg []byte) error
}

type clientConn struct {
	hm      hubpkg.HubMgr
	conn    *websocket.Conn
	sendCh  chan []byte
	closeCh chan struct{}
	msgH    mHandler
}

func (cc *clientConn) CloseConn() error {
	// TODO: clear the peer from the hub to avoid memory leak
	return cc.conn.Close()
}

func (cc *clientConn) readPump() {
	defer func() {
		cc.CloseConn()
	}()
	cc.conn.SetReadLimit(maxMessageSize)
	cc.conn.SetReadDeadline(time.Now().Add(pongWait))
	cc.conn.SetPongHandler(func(string) error {
		cc.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		msgType, msg, err := cc.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.WithError(err).Error("websocket reading failed")
			}
			break
		}

		if msgType != websocket.TextMessage {
			log.Error("only text message is allowed")
			break
		}

		if err := cc.msgH.Handle(msg); err != nil {
			log.WithError(err).Error("fail to handel message")
			return
		}
	}
}

func (cc *clientConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		cc.CloseConn()
	}()
	for {
		select {
		case message, ok := <-cc.sendCh:
			cc.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				cc.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := cc.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				log.WithError(err).Error("fail to write message")
				return
			}

			if err := w.Close(); err != nil {
				log.WithError(err).Error("fail to close message writer")
				return
			}
		case <-ticker.C:
			cc.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := cc.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithError(err).Error("fail to write ping message")
				return
			}

		case <-cc.closeCh:
			log.WithFields(log.Fields{"remote-addr": cc.conn.RemoteAddr().String()}).Info("websocket connection is closed")
			return
		}
	}
}

func (cc *clientConn) AsyncRun() {
	go cc.writePump()
	go cc.readPump()
}

func (cc *clientConn) Write(msg []byte) (int, error) {
	cc.sendCh <- msg
	return len(msg), nil
}

func (cc *clientConn) Close() error {
	cc.closeCh <- struct{}{}
	return nil
}

// NewAsyncHandler makes a new client websocekt conn
func NewAsyncHandler(conn *websocket.Conn, hm hubpkg.HubMgr) AsyncHandler {
	cc := &clientConn{
		conn:    conn,
		hm:      hm,
		sendCh:  make(chan []byte),
		closeCh: make(chan struct{}),
	}

	msgH := &msgHandler{
		hub:    nil,
		hm:     hm,
		sender: cc,
	}
	cc.msgH = msgH
	return cc
}
