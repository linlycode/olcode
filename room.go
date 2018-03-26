package olcode

type roomID string

const invalidRoomID roomID = "'"

type room struct {
	id      roomID
	creator *User
	editing *Editing
}

func (r *room) attend(user *User) {
	r.editing.Attend(user)
}

func (r *room) leave(user *User) {
	r.editing.Leave(user)
}

func (r *room) getUserList() []*User {
	return r.editing.GetUserList()
}

func (r *room) moveCursor(offset int, user *User) error {
	return r.editing.MoveCursor(offset, user)
}

func (r *room) getDocDetail() (string, map[int64]int) {
	return r.editing.GetDetail()
}

func (r *room) insertText(text string, user *User) error {
	return r.editing.Insert(text, user)
}

func (r *room) deleteText(n int, before bool, user *User) error {
	return r.editing.Delete(n, before, user)
}
