package main

import (
	"time"

	"github.com/swxctx/plex/client"
	"github.com/swxctx/plex/plog"
)

// 业务逻辑服务器
func main() {
	client.Start(&client.Config{
		InnerServers: []string{"127.0.0.1:9578", "127.0.0.1:9579"},
		ShowTrace:    true,
	})
	plog.Infof("client started...")

	for {
		time.Sleep(time.Duration(1) * time.Second)
		client.Send(&client.SendMessageArgs{
			Uid:  "1",
			Body: "{\"logic\": 1}",
			Uri:  "/logic/test/1",
		})

		time.Sleep(time.Duration(1) * time.Second)
		client.Send(&client.SendMessageArgs{
			Uid:  "2",
			Body: "{\"logic\": 2}",
			Uri:  "/logic/test/2",
		})
	}
}
