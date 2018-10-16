package jsonrpc

import "encoding/json"

// Protocol defines the struct of jsonrpc
type Protocol struct {
	Method string
	Data   interface{}
}

// Encode encodes a jsonrpc message into bytes
func Encode(method string, data interface{}) ([]byte, error) {
	p := &Protocol{
		Method: method,
		Data:   data,
	}
	return json.Marshal(p)
}

// Decode decodes a jsonrpc raw message into struct
func Decode(rawData []byte, data interface{}) (*Protocol, error) {
	p := &Protocol{
		Data: data,
	}

	if err := json.Unmarshal(rawData, p); err != nil {
		return nil, err
	}

	return p, nil
}
