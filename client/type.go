package client

import (
	"encoding/json"
	"fmt"
)

// SendMessageArgs outer send message arg
type SendMessageArgs struct {
	// 接收的uid
	Uid string `json:"uid,omitempty"`
	// 数据体[可以放json]
	Body string `json:"body,omitempty"`
	// 标识uri[不可以传系统内置的URI]
	Uri string `json:"uri,omitempty"`
}

// marshalSendMessage
func marshalSendMessage(msg *SendMessageArgs) (string, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("marshal send message err-> %v", err)
	}
	return string(data), nil
}
