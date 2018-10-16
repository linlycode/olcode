package jsonrpc

import (
	"testing"
)

type testS struct {
	UserID   int64
	UserName string
}

func TestEncodeDecode(t *testing.T) {
	method := "register"
	ts := &testS{
		UserID:   1,
		UserName: "john",
	}

	bs, err := Encode(method, ts)
	if err != nil {
		t.Fatalf("fail to encode, err=%v", err)
	}

	dTs := &testS{}
	p, err := Decode(bs, dTs)
	if err != nil {
		t.Fatalf("fail to decode, err=%v", err)
	}

	if p.Method != method {
		t.Errorf("wrong method, expect %s, got %s", method, p.Method)
	}

	if dTs.UserID != ts.UserID || dTs.UserName != ts.UserName {
		t.Errorf("wrong data, expect %v, got %v", ts, dTs)
	}
}
