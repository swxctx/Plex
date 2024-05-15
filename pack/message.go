package pack

import (
	"encoding/json"
)

// Message
type Message struct {
	// 消息序号
	Seq int64 `json:"seq,omitempty"`
	// 标识
	URI string `json:"uri,omitempty"`
	// 数据体
	Body string `json:"body,omitempty"`
}

// marshalMessage
func marshalMessage(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// unmarshalMessage
func unmarshalMessage(data []byte) (*Message, error) {
	var (
		message *Message
	)
	if err := json.Unmarshal(data, &message); err != nil {
		return nil, err
	}
	return message, nil
}
