package olcode

type roomID string

const invalidRoomID roomID = "'"

type room struct {
	id       roomID
	creator  *User
	editting *Editting
}

func (r *room) attend(user *User) {
	r.editting.Attend(user)
}

func (r *room) leave(user *User) {
	r.editting.Leave(user)
}

func (r *room) getUserList() []*UserEditting {
	return r.editting.GetUserEditingList()
}
