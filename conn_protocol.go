package olcode

import (
	"encoding/json"
)

const (
	// send
	msgSendStart = iota
	msgDocDetail
	msgRoomDeleted
	msgUserList
	msgSendEnd
)

const msgTypeLen = 4

const (
	// recv
	msgRecvStart = iota + 100
	msgDocInsert
	msgDocDelete
	msgMoveCursor
	msgRecvEnd
)

type connProtocol struct {
	MsgType int    `json:"msg_type"`
	Data    string `json:"data"`
}

func (c *connProtocol) Encode() []byte {
	return nil
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
