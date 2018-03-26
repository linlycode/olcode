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

func (c *clientRoomConn) handleMoveCursor(p *connProtocol) {
	var mvMsg moveCursorMsg
	if err := p.UnmarshalTo(&mvMsg); err != nil {
		log.Printf("fail to unmarshal move cursor message, err=%v", err)
		return
	}

	if err := c.hub.room.moveCursor(mvMsg.Offset, c.user); err != nil {
		log.Printf("fail to move cursor, err=%v", err)
		return
	}

	c.hub.broadcastDocSync()
}

func (c *clientRoomConn) handleDocInsert(p *connProtocol) {}
func (c *clientRoomConn) handleDocDelete(p *connProtocol) {}

func (c *clientRoomConn) processRecvMsg(msg []byte) {
	p := &connProtocol{}
	if err := p.Decode(msg); err != nil {
		log.Printf("fail to decode msg, msg=%v, err=%v", msg, err)
		return
	}

	var h func(p *connProtocol)
	switch p.Type() {
	case msgDocInsert:
		h = c.handleDocInsert
	case msgDocDelete:
		h = c.handleDocDelete
	case msgMoveCursor:
		h = c.handleMoveCursor
	default:
		log.Printf("should be valid recv msg type, type=%v", p.Type())
		return
	}
	h(p)
}

func (c *clientRoomConn) unregister() {
	c.hub.unregisterClientRoomConn(c)
}

func (c *clientRoomConn) readPump() {
	defer func() {
		c.conn.Close()
		c.unregister()
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
				log.Printf("fail to read message error: %v", err)
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
		c.unregister()
		c.hub.broadcastUserList()
	}()

	for {
		select {
		case p, ok := <-c.sendCh:
			log.Printf("will send msg to client")
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

			data, err := p.Encode()
			if err != nil {
				log.Printf("fail to encode msg, err=%v", err)
				continue
			}

			if _, err := w.Write(data); err != nil {
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
			//TODO: daemon for sending doc detail
		}
	}
}

// buildClientRoomConn handles websocket requests from the peer.
func buildClientRoomConn(user *User, hub *Hub, w http.ResponseWriter, r *http.Request) {
	log.Printf("start to build client room conn")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("fail to upgrade connection, err=%v", err)
		return
	}
	log.Printf("client connection upgraded")
	crConn := newClientRoomConn(user, hub, conn)
	hub.registerClientRoomConn(crConn)

	go crConn.writePump()
	go crConn.readPump()

	crConn.hub.broadcastUserList()
}
