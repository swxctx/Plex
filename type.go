package plex

import (
	"encoding/json"

	"github.com/swxctx/plex/plog"
)

// sendMessageArgs outer send message arg
type sendMessageArgs struct {
	// 接收的uid
	Uid string `json:"uid,omitempty"`
	// 数据体[可以放json]
	Body string `json:"body,omitempty"`
	// 标识uri[不可以传系统内置的URI]
	Uri string `json:"uri,omitempty"`
}

// unmarshalSendMessage
func unmarshalSendMessage(data string) *sendMessageArgs {
	var (
		msg *sendMessageArgs
	)
	if err := json.Unmarshal([]byte(data), &msg); err != nil {
		plog.Errorf("unmarshal inner message err-> %v", err)
		return nil
	}
	return msg
}
