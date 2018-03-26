package olcode

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const connMsgBufferSize = 256
const maxMessageSize = sizeLimit

const (
	writeWait  = time.Duration(2) * time.Second
	pongWait   = time.Duration(10) * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  sizeLimit,
	WriteBufferSize: sizeLimit,
}

type clientRoomConn struct {
	user *User
	conn *websocket.Conn
	hub  *Hub

	sendCh chan *connProtocol
	recvCh chan *connProtocol
}

func newClientRoomConn(user *User, hub *Hub, conn *websocket.Conn) *clientRoomConn {
	return &clientRoomConn{
		user:   user,
		conn:   conn,
		hub:    hub,
		sendCh: make(chan *connProtocol, connMsgBufferSize),
		recvCh: make(chan *connProtocol, connMsgBufferSize),
	}
}

func (c *clientRoomConn) processRecvMsg(msg []byte) {
	p := &connProtocol{}
	if err := p.Decode(msg); err != nil {
		log.Printf("fail to decode msg, msg=%v, err=%v", msg, err)
		return
	}

	t := p.Type()
	if t <= msgRecvStart || t >= msgRecvEnd {
		log.Printf("not recv message, type=%v", t)
		return
	}

	switch t {
	case msgDocInsert:
	case msgDocDelete:
	case msgMoveCursor:
	default:
		log.Panic("should be valid recve msg type range")
	}
}

func (c *clientRoomConn) readPump() {
	defer func() {
		c.hub.unregisterCh <- c
		c.conn.Close()
		c.hub.broadcastUserList()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		// TODO: filter out control message
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		c.processRecvMsg(message)
	}
}

func (c *clientRoomConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		// FIXME: may broadcast twice
		c.hub.broadcastUserList()
	}()

	for {
		select {
		case p, ok := <-c.sendCh:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("fail to get writer, err=%v", err)
				return
			}

			if _, err := w.Write(p.Encode()); err != nil {
				log.Printf("fail to write message, err=%v", err)
				return
			}

			if err := w.Close(); err != nil {
				log.Printf("fail to close writer, err=%v", err)
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("fail to write ping message, err=%v", err)
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(user *User, hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	crConn := newClientRoomConn(user, hub, conn)
	hub.registerCh <- crConn

	go crConn.writePump()
	go crConn.readPump()

	crConn.hub.broadcastUserList()
}
