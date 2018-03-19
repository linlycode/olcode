package olcode

import (
	"fmt"
	"sync"
)

// userEditting contains the editting info of the user
type userEditting struct {
	user      *User
	cursorPos int
}

// Editting is the editing of a document
type Editting struct {
	doc *Document

	uMtx          sync.Mutex
	userEdittings map[int64]*userEditting
}

// NewEditting creates an editting
func NewEditting(doc *Document, user *User) *Editting {
	ues := make(map[int64]*userEditting)
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

// Attend handles user attending the editing
func (e *Editting) Attend(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()

	if ue, ok := e.userEdittings[u.ID]; ok {
		ue.user = u
	}
	e.userEdittings[u.ID] = &userEditting{user: u}
}

// Leave handles user leaving the editting
func (e *Editting) Leave(u *User) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	delete(e.userEdittings, u.ID)
}

// cursorPosValid return whether pos is a valid cursor position in the document
func (e *Editting) cursorPosValid(pos int) bool {
	return pos >= 0 && pos <= len(e.doc.Content)
}

// CursorPosition returns user's cursor position
func (e *Editting) CursorPosition(user *User) (int, error) {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	e.doc.cMtx.Lock()
	defer e.doc.cMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return 0, fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	if !e.cursorPosValid(ue.cursorPos) {
		return 0, fmt.Errorf("user cursor position %d is invalid, userID: %d", ue.cursorPos, user.ID)
	}
	return ue.cursorPos, nil
}

// MoveCursor handles the user moving the cursor
func (e *Editting) MoveCursor(pos int, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	e.doc.cMtx.Lock()
	defer e.doc.cMtx.Unlock()

	if !e.cursorPosValid(pos) {
		return fmt.Errorf("invalid position %d", pos)
	}

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	ue.cursorPos = pos
	return nil
}

// Insert inserts str after the user's cursor
func (e *Editting) Insert(str string, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	e.doc.cMtx.Lock()
	defer e.doc.cMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	if !e.cursorPosValid(ue.cursorPos) {
		return fmt.Errorf("user cursor position %d is invalid, userID: %d", ue.cursorPos, user.ID)
	}

	e.doc.Content = e.doc.Content[:ue.cursorPos] + str + e.doc.Content[:ue.cursorPos]

	// update cursors
	if len(e.doc.Content) == len(str) {
		ue.cursorPos += len(str)
	} else {
		for _, u := range e.userEdittings {
			if u.cursorPos >= ue.cursorPos {
				u.cursorPos += len(str)
			}
		}
	}
	return nil
}

// Delete deletes n bytes before or after user's cursor
func (e *Editting) Delete(n int, before bool, user *User) error {
	e.uMtx.Lock()
	defer e.uMtx.Unlock()
	e.doc.cMtx.Lock()
	defer e.doc.cMtx.Unlock()

	ue, ok := e.userEdittings[user.ID]
	if !ok {
		return fmt.Errorf("user %d did not attend the editting", user.ID)
	}
	if !e.cursorPosValid(ue.cursorPos) {
		return fmt.Errorf("user cursor position %d is invalid, userID: %d", ue.cursorPos, user.ID)
	}

	var begin, end int
	if before {
		begin = ue.cursorPos - n
		end = ue.cursorPos

		if begin < 0 {
			begin = 0
			n = ue.cursorPos
		}
		ue.cursorPos -= n // update current user's cursor

	} else {
		begin = ue.cursorPos
		end = ue.cursorPos + n

		cLen := len(e.doc.Content)
		if end > cLen {
			end = cLen
			n = cLen - ue.cursorPos
		}
	}

	e.doc.Content = e.doc.Content[0:begin] + e.doc.Content[end:]

	// update cursors for other users
	for _, u := range e.userEdittings {
		if u.cursorPos > end {
			u.cursorPos -= n
		} else if u.cursorPos > begin && u.cursorPos < end {
			u.cursorPos = begin
		}
	}
	return nil
}
