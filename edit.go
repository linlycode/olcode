package olcode

import (
	"fmt"
	"sync"
)

// UserEditting contains the editting info of the user
type UserEditting struct {
	User      *User `json:"user"`
	CursorPos int   `json:"cursor_pos"`
}

// Editting is the editing of a document
type Editting struct {
	doc *Document

	uMtx          sync.Mutex
	userEdittings map[int64]*UserEditting
}

// NewEditting creates an editting
func NewEditting(doc *Document, user *User) *Editting {
	ues := make(map[int64]*UserEditting)
	e := &Editting{
		doc:           doc,
		userEdittings: ues,
	}
	e.Attend(user)
	return e
}

// UserCount returns the number of users
func (e *Editting) UserCount() int {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	return len(e.userEdittings)
}

// GetUserEditingList returns the information of all the editing users
func (e *Editting) GetUserEditingList() []*UserEditting {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	ues := make([]*UserEditting, 0)
	for _, ue := range e.userEdittings {
		ues = append(ues, ue)
	}
	return ues
}

// Attend handles user attending the editing
func (e *Editting) Attend(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	if _, ok := e.userEdittings[u.ID]; !ok {
		e.userEdittings[u.ID] = &UserEditting{User: u}
	}
}

// Leave handles user leaving the editting
func (e *Editting) Leave(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	delete(e.userEdittings, u.ID)
}

// CursorPosition returns user's cursor position
func (e *Editting) CursorPosition(user *User) (int, error) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	e.doc.mtx.Lock()
	defer e.doc.mtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return 0, fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	return ue.CursorPos, nil
}

// MoveCursor handles the user moving the cursor
func (e *Editting) MoveCursor(pos int, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}

	if err := e.doc.CheckOffset(pos); err != nil {
		return err
	}
	ue.CursorPos = pos
	return nil
}

// Insert inserts str after the user's cursor
func (e *Editting) Insert(str string, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	n, err := e.doc.Insert(ue.CursorPos, str)
	if err != nil {
		return err
	}

	// update cursors
	if e.doc.Len() != n {
		for _, u := range e.userEdittings {
			if u.CursorPos > ue.CursorPos {
				u.CursorPos += n
			}
		}
	}
	ue.CursorPos += n
	return nil
}

// Delete deletes n bytes before or after user's cursor
func (e *Editting) Delete(n int, before bool, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}

	begin, end, err := e.doc.Delete(ue.CursorPos, n, before)
	if err != nil {
		return err
	}

	// update cursors for other users
	for _, u := range e.userEdittings {
		if u.CursorPos >= end {
			u.CursorPos -= n
		} else if u.CursorPos > begin && u.CursorPos < end {
			u.CursorPos = begin
		}
	}
	return nil
}
