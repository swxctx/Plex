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
	plex.Start(&plex.Config{
		Port:        "9579",
		HttpPort:    "9501",
		ShowTrace:   true,
		AuthTimeout: 5,
	}, authFunc)
}
