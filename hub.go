package olcode

import (
	"encoding/json"
	"log"
)

const (
	broadcastMsgCap = 2 * 1024
)

// Hub is a hub for managing all the connections in some room
type Hub struct {
	room *room
	// map[userID]conn
	m            map[int64]*clientRoomConn
	broadcastCh  chan *connProtocol
	registerCh   chan *clientRoomConn
	unregisterCh chan *clientRoomConn
	destroyCh    chan struct{}
}

func newHub(room *room) *Hub {
	return &Hub{
		room:         room,
		m:            make(map[int64]*clientRoomConn),
		broadcastCh:  make(chan *connProtocol, broadcastMsgCap),
		registerCh:   make(chan *clientRoomConn),
		unregisterCh: make(chan *clientRoomConn),
	}
}

func (h *Hub) getRoomID() roomID {
	return h.room.id
}

func (h *Hub) registerClientRoomConn(conn *clientRoomConn) {
	user := conn.user
	if _, exist := h.m[user.ID]; exist {
		log.Printf("client room conn has existed already, roomID=%v, userID=%v", h.room.id, user.ID)
		return
	}

	if err := h.room.attend(user); err != nil {
		log.Printf("user(%v) fail to attend the room, err=%v", user.ID, err)
		return
	}

	h.m[user.ID] = conn

}

func (h *Hub) unregisterClientRoomConn(conn *clientRoomConn) {
	// TODO: need implement
}

func (h *Hub) broadcastMsg(msg *connProtocol) {
}

func (h *Hub) broadcastUserList() {
	listData, err := json.Marshal(h.room.getUserList())
	if err != nil {
		log.Printf("fail to marshal user list, err=%v", err)
		return
	}

	h.broadcastCh <- &connProtocol{
		MsgType: msgUserList,
		Data:    string(listData),
	}
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.registerCh:
			h.registerClientRoomConn(conn)
		case conn := <-h.unregisterCh:
			h.unregisterClientRoomConn(conn)
		case msg := <-h.broadcastCh:
			h.broadcastMsg(msg)
		case <-h.destroyCh:
			log.Printf("hub(%v) is detroyed", h.getRoomID())
			break
		}
	}
}
