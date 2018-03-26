package olcode

import (
	"encoding/json"
)

const (
	// send
	msgDocDetail = iota
	msgRoomDeleted
	msgUserList
)

const (
	// recv
	msgDocInsert = iota
	msgDocDelete
	msgMoveCursor
)

const msgTypeLen = 4

type connProtocol struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
}

func (c *connProtocol) Encode() ([]byte, error) {
	return json.Marshal(c)
}

func (c *connProtocol) Decode(bs []byte) error {
	return json.Unmarshal(bs, c)
}

func (c *connProtocol) Type() int {
	return c.MsgType
}

func (c *connProtocol) UnmarshalTo(v interface{}) error {
	return json.Unmarshal([]byte(c.Data), v)
}

type moveCursorMsg struct {
	Offset int `json:"offset"`
}

type docDetailMsg struct {
	Content      string         `json:"content"`
	UserEditings []*UserEditing `json:"user_editings"`
}
