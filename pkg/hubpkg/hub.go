package hubpkg

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Hub is a collection of peers(now only support two peers)
type Hub struct {
	ID   int64
	code string
	cMu  *sync.RWMutex

	pm  map[int64]*Peer
	pMu *sync.RWMutex
}

// Broadcast send the message to all the members in the hub
func (h *Hub) Broadcast(senderID int64, msg []byte) error {
	h.pMu.RLock()
	defer h.pMu.RUnlock()

	var rErr error
	for id, p := range h.pm {
		if id == senderID {
			continue
		}
		if err := p.Send(msg); err != nil {
			rErr = err
			log.WithError(err).Errorf("fail to send msg from peer#%d -> peer#%d", senderID, id)
		}
	}

	return rErr
}

// AddPeer add Peer into the hub
func (h *Hub) AddPeer(p *Peer) error {
	h.pMu.Lock()
	defer h.pMu.Unlock()

	if len(h.pm) >= 2 {
		return fmt.Errorf("hub is full")
	}
	// just override if p has existed
	h.pm[p.ID] = p
	return nil
}

func newHub(id int64) *Hub {
	return &Hub{
		ID:   id,
		code: "",
		cMu:  &sync.RWMutex{},
		pm:   make(map[int64]*Peer),
		pMu:  &sync.RWMutex{},
	}
}

// SetCode updates the code field thread-safely
func (h *Hub) SetCode(code string) error {
	h.cMu.Lock()
	defer h.cMu.Unlock()

	h.code = code
	return nil
}

// GetCode get the code field thread-safely
func (h *Hub) GetCode() string {
	h.cMu.RLock()
	defer h.cMu.RUnlock()

	return h.code
}
