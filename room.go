package olcode

type roomID string

const invalidRoomID roomID = "'"

type room struct {
	id       roomID
	creator  *User
	editting *Editting
}

func (r *room) attend(user *User) error {
	r.editting.Attend(user)
	return nil
}

func (r *room) leave(user *User) error {
	r.editting.Leave(user)
	return nil
}

func (r *room) getUserList() []*UserEditting {
	return r.editting.GetUserEditingList()
}
