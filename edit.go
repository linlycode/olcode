package olcode

import (
	"fmt"
	"sync"
)

// UserEditing contains the editting info of the user
type UserEditing struct {
	User      *User `json:"user"`
	CursorPos int   `json:"cursor_pos"`
}

// Editing is the editing of a document
type Editing struct {
	doc *Document

	uMtx          sync.RWMutex
	userEdittings map[int64]*UserEditing
}

// NewEditing creates an editting
func NewEditing(doc *Document, user *User) *Editing {
	ues := make(map[int64]*UserEditing)
	e := &Editing{
		doc:           doc,
		userEdittings: ues,
	}
	e.Attend(user)
	return e
}

// UserCount returns the number of users
func (e *Editing) UserCount() int {
	e.uMtx.RLock()
	defer e.uMtx.RUnlock()
	return len(e.userEdittings)
}

// GetUserEditingList returns the information of all the editing users
func (e *Editing) GetUserEditingList() []*UserEditing {
	e.uMtx.RLock()
	defer e.uMtx.RUnlock()
	ues := make([]*UserEditing, 0)
	for _, ue := range e.userEdittings {
		ues = append(ues, ue)
	}
	return ues
}

// Attend handles user attending the editing
func (e *Editing) Attend(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	e.userEdittings[u.ID] = &UserEditing{User: u}
}

// Leave handles user leaving the editting
func (e *Editing) Leave(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	delete(e.userEdittings, u.ID)
}

// CursorPosition returns user's cursor position
func (e *Editing) CursorPosition(user *User) (int, error) {
	e.uMtx.RLock()
	defer e.uMtx.RUnlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return 0, fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	return ue.CursorPos, nil
}

// MoveCursor handles the user moving the cursor
func (e *Editing) MoveCursor(pos int, user *User) error {
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
func (e *Editing) Insert(str string, user *User) error {
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
func (e *Editing) Delete(n int, before bool, user *User) error {
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
