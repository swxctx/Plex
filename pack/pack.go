package pack

import (
	"bytes"
	"encoding/binary"
	"github.com/swxctx/plex/plog"
	"io"
	"net"
)

const (
	// 包头长度
	headSize = 4
)

// Pack tcp write message
func Pack(message *Message) ([]byte, error) {
	// marshalMessage
	data, err := marshalMessage(message)
	if err != nil {
		return nil, err
	}

	// message size
	length := uint32(len(data))

	var (
		buffer bytes.Buffer
	)

	// write head size
	if err := binary.Write(&buffer, binary.BigEndian, length); err != nil {
		return nil, err
	}

	// write message
	_, err = buffer.Write(data)
	if err != nil {
		return nil, err
	}

	// 返回缓冲区中的字节数据
	return buffer.Bytes(), nil
}

// Unpack tcp data
func Unpack(conn net.Conn) (*Message, error) {
	// data head
	head := make([]byte, headSize)
	if _, err := io.ReadFull(conn, head); err != nil {
		if err == io.EOF {
			return nil, nil
		}
		plog.Errorf("unpack: read msg head err-> %v", err)
		return nil, err
	}

	// unpack real data length
	msgLength, err := unpackHead(head)
	if err != nil {
		plog.Errorf("unpack: unpack head err-> %v", err)
		return nil, err
	}

	if msgLength == 0 {
		plog.Warnf("unpack: msg length is zero.")
		return nil, nil
	}

	msgTemp := make([]byte, msgLength)
	if _, err := io.ReadFull(conn, msgTemp); err != nil {
		plog.Errorf("unpack: message read err-> %v", err)
		return nil, err
	}

	// real message unpack
	message, err := unmarshalMessage(msgTemp)
	if err != nil {
		plog.Errorf("unPack: unmarshalMessage err-> %v", err)
		return nil, err
	}
	return message, nil
}

// unpackHead
func unpackHead(data []byte) (uint32, error) {
	var (
		realLength uint32
	)
	// read real length
	r := bytes.NewReader(data)
	if err := binary.Read(r, binary.BigEndian, &realLength); err != nil {
		return 0, err
	}
	return realLength, nil
}
