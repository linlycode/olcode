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

func (r *room) getUserList() []*UserEditing {
	return r.editing.GetUserEditingList()
}

func (r *room) moveCursor(offset int, user *User) error {
	return r.editing.MoveCursor(offset, user)
}

func (r *room) getDocDetail() *docDetailMsg {
	content, ues := r.editing.GetDetail()
	return &docDetailMsg{Content: content, UserEditings: ues}
}
