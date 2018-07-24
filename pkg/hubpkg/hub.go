package hubpkg

import (
	"fmt"
	"sync"

	log "github.com/Sirupsen/logrus"
)

// Hub is a collection of peers(now only support two peers)
type Hub struct {
	ID int64

	pm  map[int64]*Peer
	mtx *sync.RWMutex
}

// Broadcast send the message to all the members in the hub
func (h *Hub) Broadcast(senderID int64, msg []byte) error {
	h.mtx.RLock()
	defer h.mtx.RUnlock()

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
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if len(h.pm) >= 2 {
		return fmt.Errorf("hub is full")
	}
	// just override if p has existed
	h.pm[p.ID] = p
	return nil
}

func newHub(id int64) *Hub {
	return &Hub{
		ID:  id,
		pm:  make(map[int64]*Peer),
		mtx: &sync.RWMutex{},
	}
}
