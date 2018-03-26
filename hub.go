package olcode

import (
	"encoding/json"
	"fmt"
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
	unregisterCh chan *clientRoomConn
	destroyCh    chan struct{}
}

func newHub(room *room) *Hub {
	return &Hub{
		room:        room,
		m:           make(map[int64]*clientRoomConn),
		broadcastCh: make(chan *connProtocol, broadcastMsgCap),
	}
}

func (h *Hub) getRoomID() roomID {
	return h.room.id
}

func (h *Hub) registerClientRoomConn(conn *clientRoomConn) error {
	log.Printf("register client room conn")
	user := conn.user
	if _, exist := h.m[user.ID]; exist {
		log.Printf("client room conn has existed already, roomID=%v, userID=%v", h.room.id, user.ID)
		return nil
	}

	if err := h.room.attend(user); err != nil {
		return fmt.Errorf("user(%v) fail to attend the room, err=%v", user.ID, err)
	}

	h.m[user.ID] = conn
	return nil
}

func (h *Hub) unregisterClientRoomConn(conn *clientRoomConn) {
	// TODO: need implement
}

func (h *Hub) broadcastMsg(msg *connProtocol) {
	log.Printf("broadcast message, type=%v", msg.Type())
	for _, crConn := range h.m {
		crConn.sendCh <- msg
	}
}

func (h *Hub) broadcastUserList() {
	log.Printf("broadcast user list")
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
		case msg := <-h.broadcastCh:
			h.broadcastMsg(msg)
		case <-h.destroyCh:
			log.Printf("hub(%v) is detroyed", h.getRoomID())
			break
		}
	}
}
