package main

import (
	"fmt"
	"strconv"

	"github.com/swxctx/plex"
	"github.com/swxctx/plex/plog"
)

func main() {
	authFunc := func(body string) (bool, string) {
		plog.Infof("auth, body-> %s", body)
		uid, err := strconv.ParseInt(body, 10, 64)
		if err != nil {
			plog.Errorf("ParseInt: err-> %v, str-> %s", err, body)
			return false, ""
		}
		return true, fmt.Sprintf("%d", uid)
	}

	// new server
	plex.NewServer(&plex.Config{
		Port:         "9578",
		HttpPort:     "9500",
		OuterServers: []string{"117.50.198.225:9578"},
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
