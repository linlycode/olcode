package olcode

import (
	"fmt"
	"sync"
)

type roomID string

type room struct {
	id       roomID
	creator  *User
	editting *Editting
}

type roomManager struct {
	mtx   sync.Mutex
	rooms map[roomID]*room
}

func newRoomManager() *roomManager {
	return &roomManager{
		rooms: make(map[roomID]*room),
	}
}

func (m *roomManager) create(user *User) roomID {
	var id roomID
	for {
		id = roomID(genRandString(16))
		if _, ok := m.rooms[id]; !ok {
			break
		}
	}
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.rooms[id] = &room{
		id:       id,
		creator:  user,
		editting: NewEditting(&Document{}, user),
	}
	return id
}

func (m *roomManager) attend(id roomID, user *User) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	room, ok := m.rooms[id]
	if !ok {
		return fmt.Errorf("room %v not exist", id)
	}
	room.editting.Attend(user)
	return nil
}

func (m *roomManager) leave(id roomID, user *User) error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	room, ok := m.rooms[id]
	if !ok {
		return fmt.Errorf("room %v not exist", id)
	}

	room.editting.Leave(user)

	if room.editting.UserCount() == 0 {
		delete(m.rooms, id)
	}
	return nil
}
