package olcode

import (
	"sync"
)

// User is a user
type User struct {
	ID   int64
	Name string
}

type userStore struct {
	mtx   sync.Mutex
	users map[int64]*User
}

func newUserStore() *userStore {
	return &userStore{users: make(map[int64]*User)}
}

func (s *userStore) newUser(name string) *User {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	user := &User{ID: int64(len(s.users)), Name: name}
	s.users[user.ID] = user
	return user
}

func (s *userStore) get(userID int64) *User {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	user, ok := s.users[userID]
	if !ok {
		return nil
	}
	return user
}
