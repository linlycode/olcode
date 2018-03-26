package olcode

import (
	"encoding/json"
	"log"
	"sync"
)

const (
	broadcastMsgCap = 2 * 1024
)

// Hub is a hub for managing all the connections in some room
type Hub struct {
	room *room
	// map[userID]conn
	m           map[int64]*clientRoomConn
	mMtx        sync.RWMutex
	broadcastCh chan *connProtocol
	destroyCh   chan struct{}
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

func (h *Hub) registerClientRoomConn(conn *clientRoomConn) {
	log.Printf("register client room conn, userID=%v, roomID=%v", conn.user.ID, h.room.id)
	h.mMtx.Lock()
	defer h.mMtx.Unlock()
	if _, exist := h.m[conn.user.ID]; exist {
		log.Printf("client room conn has existed already, roomID=%v, userID=%v", h.room.id, conn.user.ID)
		return
	}
	h.room.attend(conn.user)
	h.m[conn.user.ID] = conn
}

func (h *Hub) unregisterClientRoomConn(conn *clientRoomConn) {
	log.Printf("unregister client room conn, userID=%v, roomID=%v", conn.user.ID, h.room.id)
	h.mMtx.Lock()
	defer h.mMtx.Unlock()

	h.room.leave(conn.user)
	delete(h.m, conn.user.ID)
}

func (h *Hub) broadcastMsg(msg *connProtocol) {
	log.Printf("broadcast message, type=%v", msg.Type())
	h.mMtx.RLock()
	defer h.mMtx.RUnlock()
	for _, crConn := range h.m {
		crConn.sendCh <- msg
	}
}

// TODO: broadcast to everyone excluding sender
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

// TODO: broadcast to everyone excluding sender
func (h *Hub) broadcastDocDetail() {
	log.Printf("broadcast doc detail")
	docDetail, err := json.Marshal(h.room.getDocDetail())
	if err != nil {
		log.Printf("fail to marsh doc detail, err=%v", err)
		return
	}

	h.broadcastCh <- &connProtocol{
		MsgType: msgDocDetail,
		Data:    string(docDetail),
	}
}

func (h *Hub) stopRun() {
	close(h.destroyCh)
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
