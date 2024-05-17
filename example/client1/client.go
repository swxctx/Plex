package main

import (
	"io"
	"net"
	"time"

	"github.com/swxctx/plex/pack"
	"github.com/swxctx/plex/plog"
)

func main() {
	// connect plex server
	conn, err := net.Dial("tcp", "127.0.0.1:9579")
	if err != nil {
		plog.Errorf("err-> %v", err)
		return
	}
	plog.Infof("connected.")

	authSuccess := false
	// 循环接收消息
	for {
		if !authSuccess {
			// 发送消息
			writeData, err := pack.Pack(&pack.Message{
				URI:  "/auth/server",
				Body: "plex-example-1",
			})
			if err != nil {
				plog.Errorf("pack err-> %v", err)
				return
			}
			conn.Write(writeData)
		}

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
		if message.URI == "/auth/success" {
			authSuccess = true
			go startClientHeartbeat(conn)
			plog.Infof("auth success")
		}
	}
}

// startClientHeartbeat
func startClientHeartbeat(conn net.Conn) {
	msg, err := pack.Pack(&pack.Message{
		URI: "/heartbeat",
	})
	if err != nil {
		plog.Errorf("client heartbeat err-> %v", err)
		return
	}

	// time ticket
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if _, err := conn.Write(msg); err != nil {
				plog.Errorf("client heartbeat send err-> %v", err)
				return
			}
		}
	}
}
