package olcode

import (
	"errors"
	"fmt"
	"sync"
)

var errRoomNotExist = errors.New("room not exist")

// HubMgr is the manager for all the hubs
type HubMgr struct {
	mtx  sync.Mutex
	hubs map[roomID]*Hub
}

// NewHubMgr is the build function of HubMgr
func NewHubMgr() *HubMgr {
	return &HubMgr{
		hubs: make(map[roomID]*Hub),
	}
}

func (m *HubMgr) registerHub(user *User) (roomID, error) {
	var id roomID
	for {
		id = roomID(genRandString(16))
		if _, ok := m.hubs[id]; !ok {
			break
		}
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()

	r := &room{
		id:      id,
		creator: user,
		editing: NewEditing(&Document{}, user),
	}

	h := newHub(r)
	m.hubs[id] = h

	go h.run()
	return id, nil
}

func (m *HubMgr) unregisterHub(id roomID, user *User) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	hub, ok := m.hubs[id]
	if !ok {
		return fmt.Errorf("hub %v not exist", id)
	}

	if user.ID != hub.room.creator.ID {
		return fmt.Errorf("non-creator can not delete the room, operatorID=%v", user.ID)
	}

	delete(m.hubs, id)
	return nil
}

func (m *HubMgr) getHub(id roomID) (*Hub, error) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	hub, exist := m.hubs[id]
	if !exist {
		return nil, errRoomNotExist
	}
	return hub, nil
}
