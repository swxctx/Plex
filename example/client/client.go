package main

import (
	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
	"io"
	"net"
	"time"
)

func main() {
	// connect plex server
	conn, err := net.Dial("tcp", "127.0.0.1:9578")
	if err != nil {
		plog.Errorf("err-> %v", err)
		return
	}
	plog.Infof("connected.")

	//i := 0
	// 循环接收消息
	for {
		//// 发送消息
		//writeData, err := pack.Pack(&pack.Message{
		//	Seq:  1,
		//	URI:  fmt.Sprintf("/auth/request_%d", i),
		//	Body: fmt.Sprintf("r: %d", i),
		//})
		//if err != nil {
		//	plog.Errorf("pack err-> %v", err)
		//	return
		//}
		//conn.Write(writeData)
		//i++

		message, err := pack.Unpack(conn)
		if err != nil {
			plog.Errorf("unpack err-> %v", err)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				conn.Close()
				break
			}
			continue
		}
		if message == nil {
			continue
		}
		plog.Infof("receive message-> %#v", message)

		time.Sleep(time.Duration(1) * time.Second)
	}
}
