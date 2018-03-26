package olcode

type roomID string

const invalidRoomID roomID = "'"

type room struct {
	id       roomID
	creator  *User
	editting *Editing
}

func (r *room) attend(user *User) {
	r.editting.Attend(user)
}

func (r *room) leave(user *User) {
	r.editting.Leave(user)
}

func (r *room) getUserList() []*UserEditing {
	return r.editting.GetUserEditingList()
}

func (r *room) moveCursor(offset int, user *User) error {
	return r.editting.MoveCursor(offset, user)
}
