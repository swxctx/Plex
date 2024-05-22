package main

import (
	"github.com/swxctx/plex"
	"github.com/swxctx/plex/plog"
)

func main() {
	authFunc := func(body string) (bool, string) {
		plog.Infof("auth, body-> %s", body)
		if body == "plex-example" {
			return true, "1"
		}
		if body == "plex-example-1" {
			return true, "2"
		}
		return false, ""
	}

	// new server
	plex.NewServer(&plex.Config{
		Port:         "9578",
		HttpPort:     "9500",
		OuterServers: []string{"112.56.77.90:9578", "112.56.77.91:9579"},
		ShowTrace:    true,
		AuthTimeout:  5,
	}, authFunc)

	// 上下线订阅
	plex.SetOnlineStatusSubscribe(func(isOnline bool, uid string) {
		plog.Infof("上线状态: %v, UID-> %s", isOnline, uid)
	})

	// start
	plex.Start()
}
